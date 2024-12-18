package receiver

import (
	"bytes"
	"context"
	"digiauth/pkg/main-app/db"
	sql "digiauth/pkg/main-app/db/sqlconfig"
	models "digiauth/pkg/main-app/user/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func GetConnections(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	var req models.GetConnectionsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	queries := sql.New(db.DB)
	connections, conerr := queries.GetConnectionsByUserID(ctx, req.Id)
	if conerr != nil {
		log.Println("Error getting connection to db : ", conerr.Error())
		http.Error(w, "Error getting connection to db : "+conerr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"connections": connections})
}

func GetCredentials(w http.ResponseWriter, r *http.Request) {

	resp, err := http.Get("http://localhost:6041/credentials")
	if err != nil {
		http.Error(w, "Failed to contact external service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the response from the external service
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	// Return the response from the external service to the original caller
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func ReceiveInvitation(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	var requestData models.ReceiveInvitationRequest
	//Decode the request body into req struct
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Convert the req struct to JSON for the external request
	requestBody, err := json.Marshal(requestData.Invitation)
	if err != nil {
		http.Error(w, "Failed to marshal request", http.StatusInternalServerError)
		return
	}

	resp, err := http.Post("http://localhost:6041/connections/receive-invitation", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		http.Error(w, "Failed to contact external service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the response from the external service
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	var responseData models.ResponseReceiveInvitation
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		http.Error(w, "Failed to parse response", http.StatusInternalServerError)
		return
	}
	log.Println("response data for receiving: ", responseData)
	queries := sql.New(db.DB)
	insertDBErr := queries.CreateConnection(ctx, sql.CreateConnectionParams{
		ConnectionID: responseData.ConnectionID,
		ID:           requestData.UserID,
		MyMailID:     requestData.MyMailId,
		TheirMailID:  requestData.TheirMailId,
	})
	if insertDBErr != nil {
		log.Println("Error inserting connection to db : ", insertDBErr.Error())
		http.Error(w, "Error inserting connection to db : "+insertDBErr.Error(), http.StatusInternalServerError)
		return

	}

	// Return the response from the external service to the original caller
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

// This is the function to create invitation for connection
func CreateInvitation(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var requestData models.CreateSendInvitationRequest

	// Decode the JSON request body into the struct
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Println("request data: ", requestData)

	resp, err := http.Post("http://localhost:8041/connections/create-invitation", "application/json", bytes.NewBuffer([]byte{}))
	if err != nil {
		http.Error(w, "Failed to contact external service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the response from the external service
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	// Parse the JSON response
	var responseData models.InvitationResponse

	err = json.Unmarshal(body, &responseData)
	if err != nil {
		http.Error(w, "Failed to parse response", http.StatusInternalServerError)
		return
	}

	// Print the JSON string
	invitationString, marshalErr := json.Marshal(responseData.Invitation)
	if marshalErr != nil {
		fmt.Println("Error:", marshalErr)
		return
	}

	queries := sql.New(db.DB)
	insertDBErr := queries.CreateConnection(ctx, sql.CreateConnectionParams{
		ConnectionID: responseData.ConnectionID,
		ID:           requestData.Id,
		MyMailID:     requestData.MyMailId,
		TheirMailID:  requestData.TheirMailId,
	})

	if insertDBErr != nil {
		log.Println("Error inserting connection to db : ", insertDBErr.Error())
		http.Error(w, "Error inserting connection to db : "+insertDBErr.Error(), http.StatusInternalServerError)
		return
	}

	emailBody := models.SendEmail{
		Email: requestData.TheirMailId,
		Message: models.Message{
			Subject: "Invitation to connect",
			Body:    "{<h5 style=\"margin:0;padding:0\">\"invitation\":" + string(invitationString) + ",<br/>\"their_mail_id\":\"" + requestData.MyMailId + "\"</h5>}",
		},
	}

	emailPayload, err := json.Marshal(emailBody)
	if err != nil {
		fmt.Println("Error in Marshalling emailbody:", err)
		return
	}

	_, mailErr := http.Post("https://q648rhgza1.execute-api.ap-south-1.amazonaws.com/prod", "application/json", bytes.NewBuffer(emailPayload))
	if mailErr != nil {
		log.Println("Failed to send email in create invitation")
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Invitation Sent Successfully"}`))
}

// This is the function for registering DID with Ledger
func RegisterDID(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterDIDRequest

	// Decode the request body into the req struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Convert the req struct to JSON for the external request
	requestBody, err := json.Marshal(req)
	if err != nil {
		http.Error(w, "Failed to marshal request", http.StatusInternalServerError)
		return
	}

	// Send the request to the external endpoint
	resp, err := http.Post("http://test.bcovrin.vonx.io/register", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		http.Error(w, "Failed to contact external service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the response from the external service
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	// Return the response from the external service to the original caller
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func GetRecords(ConnectionId string) (models.ProofRecord, error) {
	var allRecords models.ProofRecords
	var singleRecord models.ProofRecord
	var responseBody models.ProofRecord
	resp, err := http.Get("http://localhost:6041/present-proof-2.0/records")
	if err != nil {
		log.Println("Failed to contact external service")
		return models.ProofRecord{}, err
	}
	defer resp.Body.Close()

	// Read the response from the external service
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read response")
		return models.ProofRecord{}, err
	}

	err = json.Unmarshal(body, &allRecords)
	if err != nil {
		log.Println("Failed to unmarshall body into allRecords")
		return models.ProofRecord{}, err
	}

	for i := 0; i < len(allRecords.Results); i++ {
		err = json.Unmarshal(body, &singleRecord)
		if err != nil {
			log.Println("Failed to unmarshall body into singleRecord")
			return models.ProofRecord{}, err
		}
		if singleRecord.ConnectionId == ConnectionId {
			if singleRecord.State == "request-received" {
				responseBody = singleRecord
			}
		}
	}

	return responseBody, nil
}

func SendPresentation(w http.ResponseWriter, r *http.Request) {
	var req models.SendPresentationRequest
	var sendPresentationResponse models.SendPresentationResponse

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	record, err := GetRecords(req.ConnectionID)
	if err != nil {
		log.Println("GetRecords not working properly")
		http.Error(w, "Failed to get records", http.StatusInternalServerError)
		return
	}

	url := fmt.Sprintf("http://localhost:6041/present-proof-2.0/records/%s/send-presentation", record.Pres_Ex_Id)

	requestBody, err := json.Marshal(req)
	if err != nil {
		http.Error(w, "Failed to marshal request", http.StatusInternalServerError)
		return
	}

	response, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		http.Error(w, "Failed to send presentation", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		http.Error(w, "Failed to read response from send presentation", http.StatusInternalServerError)
		return
	}

	if response.StatusCode != http.StatusOK {
		log.Printf("Non-200 Response: %s, Body: %s", response.Status, string(body))
		http.Error(w, "Failed to send presentation", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &sendPresentationResponse)
	if err != nil {
		http.Error(w, "Failed to unmarshal body from send presentation", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"response": sendPresentationResponse})
}

package verifier

import (
	"bytes"
	"context"
	"digiauth/pkg/main-app/db"
	sql "digiauth/pkg/main-app/db/sqlconfig"
	models "digiauth/pkg/main-app/verifier/models"
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

func SendProofRequest(w http.ResponseWriter, r *http.Request) {
	var req models.SendProofRequestRequest

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

	resp, err := http.Post("http://localhost:4041/present-proof-2.0/send-request", "application/json", bytes.NewBuffer(requestBody))
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

func GetSchemasDB(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Initialize the SQL queries struct
	queries := sql.New(db.DB)

	// Fetch the schema by ID from the database using GetSchemaById
	res, err := queries.GetSchema(ctx)
	if err != nil {
		log.Println("Error fetching schema from db:", err.Error())
		http.Error(w, "Error fetching schema from db: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response header and encode the schema response as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"schema": res})
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

// ! Search mymailId and theirmailid in connections
// ! Make requests to /records(using their_mail_id)
// ! check connection state to be "done" but only check for the matching connection id that were retrieved previously.
// ! If state is done for any of the records then return "true" or "verified" for that connection
func VerifyPresentation(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var request models.VerifyPresentationRequest
	// Decode the request body into the req struct
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	queries := sql.New(db.DB)
	connections, getDBErr := queries.FetchConnections(ctx, sql.FetchConnectionsParams{
		TheirMailID: request.MyMailId,
		MyMailID:    request.TheirMailID,
	})

	if getDBErr != nil {
		log.Println("Error inserting connection to db : ", getDBErr.Error())
		http.Error(w, "Error inserting connection to db : "+getDBErr.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := http.Get("http://localhost:6041/present-proof-2.0/records")
	if err != nil {
		log.Println("Failed to contact external service")
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read response")
		return
	}

	var responseData models.GetRecordsResponse
	marshalErr := json.Unmarshal(body, &responseData)
	if marshalErr != nil {
		log.Println("Error marshalling records response", marshalErr)
		http.Error(w, "Error marshalling records response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	for i := 0; i < len(responseData.Results); i++ {
		if responseData.Results[i].ConnectionID == connections[0].ConnectionID {
			if responseData.Results[i].State == "done" {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"message": "Presentation Verified Successfully"}`))
				return
			}
		}
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"message": "Presentation Not Verified"}`))
	defer resp.Body.Close()

}

package receiver

import (
	"bytes"
	"context"
	"digiauth/main-app/db"
	sql "digiauth/main-app/db/sqlconfig"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type ResponseReceiveInvitation struct {
	State              string `json:"state"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
	ConnectionID       string `json:"connection_id"`
	MyDID              string `json:"my_did"`
	TheirLabel         string `json:"their_label"`
	TheirRole          string `json:"their_role"`
	ConnectionProtocol string `json:"connection_protocol"`
	RFC23State         string `json:"rfc23_state"`
	InvitationKey      string `json:"invitation_key"`
	InvitationMsgID    string `json:"invitation_msg_id"`
	RequestID          string `json:"request_id"`
	Accept             string `json:"accept"`
	InvitationMode     string `json:"invitation_mode"`
}

type RegisterDIDRequest struct {
	Seed  string `json:"seed"`
	Alias string `json:"alias"`
	Role  string `json:"Role"`
}

type CreateSendInvitationRequest struct {
	Alias              string   `json:"alias"`
	HandshakeProtocols []string `json:"handshake_protocols"`
	MyLabel            string   `json:"my_label"`
}

// This is for receiving invitation services
type Service struct {
	Id              string   `json:"id"`
	Type            string   `json:"type"`
	RecipientKeys   []string `json:"recipientKeys"`
	ServiceEndpoint string   `json:"serviceEndpoint"`
}

type ReceiveInvitationRequest struct {
	Type            string   `json:"@type"`
	RecipientKeys   []string `json:"recipientKeys"`
	Id              string   `json:"@id"`
	UserID          int64    `json:"id"`
	Label           string   `json:"label"`
	ServiceEndpoint string   `json:"serviceEndpoint"`
}

type GetConnectionsRequest struct {
	Id int64 `json:"id"`
}

func GetConnections(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	var req GetConnectionsRequest
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

// This is the function to receive invitation for connection (rn status==deleted rest working fine)
func ReceiveInvitation(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	var requestData ReceiveInvitationRequest
	//Decode the request body into req struct
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Convert the req struct to JSON for the external request
	requestBody, err := json.Marshal(requestData)
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

	var responseData ResponseReceiveInvitation
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		http.Error(w, "Failed to parse response", http.StatusInternalServerError)
		return
	}

	queries := sql.New(db.DB)
	insertDBErr := queries.CreateConnection(ctx, sql.CreateConnectionParams{
		ConnectionID: responseData.ConnectionID,
		ID:           requestData.UserID,
		Alias:        responseData.TheirLabel,
		MyRole:       "invitee",
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
	var req CreateSendInvitationRequest
	//Decode the request body into req struct
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

	resp, err := http.Post("http://localhost:6041/connections/create-invitation", "application/json", bytes.NewBuffer(requestBody))
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

// This is the function for registering DID with Ledger
func RegisterDID(w http.ResponseWriter, r *http.Request) {
	var req RegisterDIDRequest

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

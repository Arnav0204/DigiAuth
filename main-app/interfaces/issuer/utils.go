package issuer

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

type RegisterDIDRequest struct {
	Seed  string `json:"seed"`
	Alias string `json:"alias"`
	Role  string `json:"Role"`
}

type RegisterSchemaRequest struct {
	Attributes    []string `json:"attributes"`
	SchemaName    string   `json:"schema_name"`
	SchemaVersion string   `json:"schema_version"`
}

type CreateCredentialDefinationRequest struct {
	Schemaid string `json:"schema_id"`
	Tag      string `json:"tag"`
}

type CreateSendInvitationRequest struct {
	Alias string `json:"alias"`
	Id    int64  `json:"id"`
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

type IndyFilter struct {
	CredDefID string `json:"cred_def_id"`
}

type Filter struct {
	Indy IndyFilter `json:"indy"`
}

type CredentialAttribute struct {
	MimeType string `json:"mime-type"`
	Name     string `json:"name"`
	Value    string `json:"value"`
}

type CredentialPreview struct {
	Type       string                `json:"@type"`
	Attributes []CredentialAttribute `json:"attributes"`
}

type CredentialIssuance struct {
	ConnectionID      string            `json:"connection_id"`
	Filter            Filter            `json:"filter"`
	CredentialPreview CredentialPreview `json:"credential_preview"`
	SchemaIssuerDID   string            `json:"schema_issuer_did"`
	SchemaVersion     string            `json:"schema_version"`
	SchemaID          string            `json:"schema_id"`
	SchemaName        string            `json:"schema_name"`
	IssuerDID         string            `json:"issuer_did"`
}

type GetConnectionsRequest struct {
	Id int64 `json:"id"`
}

func IssueCredential(w http.ResponseWriter, r *http.Request) {
	var req CredentialIssuance
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

	resp, err := http.Post("http://localhost:8041/issue-credential-2.0/send", "application/json", bytes.NewBuffer(requestBody))
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
		log.Println("Error inserting connection to db : ", conerr.Error())
		http.Error(w, "Error inserting connection to db : "+conerr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"connections": connections})
}

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

	resp, err := http.Post("http://localhost:8041/connections/receive-invitation", "application/json", bytes.NewBuffer(requestBody))
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
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var requestData CreateSendInvitationRequest

	// Decode the JSON request body into the struct
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

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
	var responseData struct {
		ConnectionID string `json:"connection_id"`
	}

	err = json.Unmarshal(body, &responseData)
	if err != nil {
		http.Error(w, "Failed to parse response", http.StatusInternalServerError)
		return
	}

	queries := sql.New(db.DB)
	insertDBErr := queries.CreateConnection(ctx, sql.CreateConnectionParams{
		ConnectionID: responseData.ConnectionID,
		ID:           requestData.Id,
		Alias:        requestData.Alias,
		MyRole:       "inviter",
	})

	if insertDBErr != nil {
		log.Println("Error inserting connection to db : ", insertDBErr.Error())
		http.Error(w, "Error inserting connection to db : "+insertDBErr.Error(), http.StatusInternalServerError)
		return

	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

// This is a function that registers schema with ledger
func RegisterSchema(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var req RegisterSchemaRequest
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

	registerSchemaResp, err := http.Post("http://localhost:8041/schemas", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		http.Error(w, "Failed to contact external service", http.StatusInternalServerError)
		return
	}
	defer registerSchemaResp.Body.Close()

	// Read the response from the external service
	registerSchemaBody, err := io.ReadAll(registerSchemaResp.Body)
	if err != nil {
		http.Error(w, "Failed to read response from register schema route", http.StatusInternalServerError)
		return
	}

	var registerSchemaResponseData struct {
		SchemaId string `json:"schema_id"`
	}
	err = json.Unmarshal(registerSchemaBody, &registerSchemaResponseData)
	if err != nil {
		http.Error(w, "Failed to parse response", http.StatusInternalServerError)
		return
	}

	createCredentialDefinationRequestBody := map[string]interface{}{
		"schema_id":                registerSchemaResponseData.SchemaId,
		"tag":                      "default",
		"support_revocation":       true,
		"revocation_registry_size": 1000,
	}

	// Marshal the map into a JSON byte slice
	RequestBody, err := json.Marshal(createCredentialDefinationRequestBody)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	createCredentialDefinationResp, err := http.Post("http://localhost:8041/credential-definitions", "application/json", bytes.NewBuffer(RequestBody))
	if err != nil {
		http.Error(w, "Failed to contact external service", http.StatusInternalServerError)
		return
	}
	defer createCredentialDefinationResp.Body.Close()

	createCredentialDefinationBody, err := io.ReadAll(createCredentialDefinationResp.Body)
	if err != nil {
		http.Error(w, "Failed to read response from create credential definition route", http.StatusInternalServerError)
		return
	}

	var createCredentialDefinationResponseData struct {
		CrendentialDefinitionId string `json:"crendential_definition_id"`
	}
	err = json.Unmarshal(createCredentialDefinationBody, &createCredentialDefinationResponseData)
	if err != nil {
		http.Error(w, "Failed to parse response", http.StatusInternalServerError)
		return
	}

	queries := sql.New(db.DB)
	insertDBErr := queries.CreateSchema(ctx, sql.CreateSchemaParams{
		SchemaID:               registerSchemaResponseData.SchemaId,
		CredentialDefinitionID: createCredentialDefinationResponseData.CrendentialDefinitionId,
		SchemaName:             req.SchemaName,
	})
	if insertDBErr != nil {
		log.Println("Error inserting connection to db : ", insertDBErr.Error())
		http.Error(w, "Error inserting connection to db : "+insertDBErr.Error(), http.StatusInternalServerError)
		return

	}

	// Return the response from the external service to the original caller
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Schema registered successfully"}`))
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

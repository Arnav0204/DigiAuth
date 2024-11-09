package issuer

import (
	"bytes"
	"context"
	"digiauth/pkg/main-app/db"
	sql "digiauth/pkg/main-app/db/sqlconfig"
	models "digiauth/pkg/main-app/issuer/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func IssueCredential(w http.ResponseWriter, r *http.Request) {
	var req models.IssueCredentialRequest
	// Decode the request body into the req struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var did string = "V8ErVjLajTWdW5CH1jKPyt"

	requestBody := models.CredentialIssuance{
		ConnectionID: req.ConnectionID,
		Filter: map[string]models.IndyFilter{
			"indy": {
				CredDefID: req.CredentialDefinitionId,
			},
		},
		CredentialPreview: models.CredentialPreview{
			Type:       "https://didcomm.org/issue-credential/2.0/credential-preview",
			Attributes: req.Attributes,
		},
		SchemaIssuerDID: did,
		SchemaID:        req.SchemaId,
		SchemaName:      req.SchemaName,
		IssuerDID:       did,
	}

	// Convert the req struct to JSON for the external request
	ledgerRequest, err := json.Marshal(requestBody)
	if err != nil {
		http.Error(w, "Failed to marshal request", http.StatusInternalServerError)
		return
	}

	resp, err := http.Post("http://localhost:8041/issue-credential-2.0/send", "application/json", bytes.NewBuffer(ledgerRequest))
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
	var req models.GetConnectionsRequest
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

func GetSchemasDB(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Define the request body model (for SchemaIdDB)
	var req models.SchemaIdDB
	// Initialize the SQL queries struct
	queries := sql.New(db.DB)

	// Decode the request body into SchemaIdDB
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Fetch the schema by ID from the database using GetSchemaById
	res, err := queries.GetSchemaById(ctx, req.Id)
	if err != nil {
		log.Println("Error fetching schema from db:", err.Error())
		http.Error(w, "Error fetching schema from db: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response header and encode the schema response as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"schema": res})
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
	// log.Println("receiv")
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

	var responseData models.ResponseReceiveInvitation
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

	fmt.Println(body)
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

	var req models.RegisterSchemaRequest
	// Decode the request body into the req struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var Tag = req.SchemaName
	log.Println(Tag)
	log.Println("req: ", req)
	// Convert the req struct to JSON for the external request
	requestBody, err := json.Marshal(req)
	if err != nil {
		http.Error(w, "Failed to marshal request", http.StatusInternalServerError)
		return
	}

	// log.Println("req schema after marshal: ",requestBody)
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

	log.Println("resp: ", registerSchemaResp)
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
		"tag":                      Tag,
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
		CrendentialDefinitionId string `json:"credential_definition_id"`
	}
	log.Println(createCredentialDefinationResponseData.CrendentialDefinitionId)

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
		Attributes:             req.Attributes,
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

func GetSchemas(w http.ResponseWriter, r *http.Request) {
	// Make the GET request to the external endpoint to fetch schemas
	resp, err := http.Get("http://localhost:8041/schemas/created")
	if err != nil {
		http.Error(w, "Failed to fetch schemas from external service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to fetch schemas: unexpected status code", http.StatusInternalServerError)
		return
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	// Define a struct to capture the schema IDs from the response
	var schemasResponse struct {
		SchemaIds []string `json:"schema_ids"`
	}

	// Parse the JSON response into the struct
	err = json.Unmarshal(body, &schemasResponse)
	if err != nil {
		http.Error(w, "Failed to parse schemas response", http.StatusInternalServerError)
		return
	}

	// Log the received schema IDs for debugging
	log.Printf("Fetched schema IDs: %v\n", schemasResponse.SchemaIds)

	// Send the schema IDs back as a JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"schema_ids": schemasResponse.SchemaIds})
}

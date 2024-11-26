package issuer

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
	Id          int64  `json:"id"`
	MyMailId    string `json:"my_mail_id"`
	TheirMailId string `json:"their_mail_id"`
}

type Invitation struct {
	Type            string   `json:"@type"`
	ID              string   `json:"@id"`
	Label           string   `json:"label"`
	RecipientKeys   []string `json:"recipientKeys"`
	ServiceEndpoint string   `json:"serviceEndpoint"`
}

type InvitationResponse struct {
	ConnectionID string     `json:"connection_id"`
	Invitation   Invitation `json:"invitation"`
}

type SchemaIdDB struct {
	Id string `json:"id"`
}

type SchemaResponse struct {
	SchemaID               string   `json:"schema_id"`
	CredentialDefinitionID string   `json:"credential_definition_id"`
	SchemaName             string   `json:"schema_name"`
	Attributes             []string `json:"attributes"`
}

// This is for receiving invitation services
type Service struct {
	Id              string   `json:"id"`
	Type            string   `json:"type"`
	RecipientKeys   []string `json:"recipientKeys"`
	ServiceEndpoint string   `json:"serviceEndpoint"`
}

type ReceiveInvitationRequest struct {
	UserID      int64      `json:"id"`
	TheirMailId string     `json:"their_mail_id"`
	MyMailId    string     `json:"my_mail_id"`
	Invitation  Invitation `json:"invitation"`
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

type GetConnectionsRequest struct {
	Id int64 `json:"id"`
}

type IssueCredentialRequest struct {
	ConnectionID           string                `json:"connection_id"`
	SchemaName             string                `json:"schema_name"`
	SchemaId               string                `json:"schema_id"`
	CredentialDefinitionId string                `json:"credential_definition_id"`
	Attributes             []CredentialAttribute `json:"attributes"`
}

type CredentialIssuance struct {
	ConnectionID      string                `json:"connection_id"`
	Filter            map[string]IndyFilter `json:"filter"`
	CredentialPreview CredentialPreview     `json:"credential_preview"`
	SchemaIssuerDID   string                `json:"schema_issuer_did"`
	SchemaID          string                `json:"schema_id"`
	SchemaName        string                `json:"schema_name"`
	IssuerDID         string                `json:"issuer_did"`
}

type SendEmail struct {
	Email   string  `json:"email"`
	Message Message `json:"message"`
}

type Message struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

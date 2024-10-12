package verifier

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

type GetConnectionsRequest struct {
	Id int64 `json:"id"`
}

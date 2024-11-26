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
	Id          int64  `json:"id"`
	MyMailId    string `json:"my_mail_id"`
	TheirMailId string `json:"their_mail_id"`
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

type GetConnectionsRequest struct {
	Id int64 `json:"id"`
}

type Restriction struct {
	CredDefID string `json:"cred_def_id"`
}

type RequestedAttribute struct {
	Name         string        `json:"name"`
	Restrictions []Restriction `json:"restrictions"`
}

type IndyReq struct {
	Name                string                        `json:"name"`
	Version             string                        `json:"version"`
	RequestedAttributes map[string]RequestedAttribute `json:"requested_attributes"`
	RequestedPredicates []interface{}                 `json:"requested_predicates"`
}

type PresentationReq struct {
	Indy IndyReq `json:"indy"`
}

type SendProofRequestRequest struct {
	ConnectionID        string          `json:"connection_id"`
	PresentationRequest PresentationReq `json:"presentation_request"`
	Trace               bool            `json:"trace"`
}

type ProofRecord struct {
	Pres_Ex_Id   string `json:"pres_ex_id"`
	State        string `json:"state"`
	ConnectionId string `json:"connection_id"`
}

type ProofRecords struct {
	Results []ProofRecord `json:"results"`
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
type SendEmail struct {
	Email   string  `json:"email"`
	Message Message `json:"message"`
}

type Message struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

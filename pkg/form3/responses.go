package form3

//StatusAble Interface
type StatusAble interface {
	SetStatus(status int)
}

// BaseHTTPResponse struct
type BaseHTTPResponse struct {
	StatusCode int `json:"statusCode,omitempty"`
}

// SetStatus receiver on BaseHTTPResponse object
func (b *BaseHTTPResponse) SetStatus(status int) {
	b.StatusCode = status
}

// ListResponse object
type ListResponse struct {
	BaseHTTPResponse
	Data []Account `json:"data"`
}

// SetStatus receiver on ListResponse
func (b *ListResponse) SetStatus(status int) {
	b.StatusCode = status
}

// Account data struct
type Account struct {
	Type           string     `json:"type"`
	ID             string     `json:"id"`
	OrganisationID string     `json:"organisation_id"`
	Version        int        `json:"version"`
	Attributes     Attributes `json:"attributes,omitempty"`
}

type CreateAccount struct {
	Type           string     `json:"type"`
	ID             string     `json:"id"`
	OrganisationID string     `json:"organisation_id"`
	Attributes     CreateAccountAttributes  `json:"attributes"`
}

type CreateAccountAttributes struct {
	Country                     string   `json:"country"`
	BaseCurrency                string   `json:"base_currency"`
	BankID                      int      `json:"bank_id"`
	BankIDCode                  string   `json:"bank_id_code"`
	Bic                         string   `json:"bic"`
}

//Attributes struct for the returned data
type Attributes struct {
	Country                     string   `json:"country"`
	BaseCurrency                string   `json:"base_currency"`
	AccountNumber               string   `json:"account_number"`
	BankID                      int      `json:"bank_id"`
	BankIDCode                  string   `json:"bank_id_code"`
	Bic                         string   `json:"bic"`
	Iban                        string   `json:"iban"`
	AccountClassification       string   `json:"account_classification,omitempty"`
	JointAccount                bool     `json:"joint_account,omitempty"`
	AccountMatchingOptOut       bool     `json:"account_matching_opt_out,omitempty"`
	Title                       string   `json:"title,omitempty"`
	FirstName                   string   `json:"first_name,omitempty"`
	AlternativeBankAccountNames []string `json:"alternative_bank_account_names,omitempty"`
	SecondaryIdentification     string   `json:"secondary_identification,omitempty"`
}

// Error struct
type Error struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// Errors struct
type Errors struct {
	FieldErrors   map[string][]Error `json:"fieldErrors,omitempty"`
	GeneralErrors []Error            `json:"generalErrors,omitempty"`
}

//FetchAccountResponse struct
type FetchAccountResponse struct {
	BaseHTTPResponse
	Data Account `json:"data"`
}

// CreateAccountRequest struct
type CreateAccountRequest struct {
	Data CreateAccount `json:"data"`
}

// CreateAccountResponse struct
type CreateAccountResponse struct {
	BaseHTTPResponse
	Data Account `json:"data"`
}
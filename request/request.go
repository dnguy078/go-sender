package request

type EmailRequest struct {
	ToEmail   string `json:"toEmail"`
	FromEmail string `json:"fromEmail"`
	Subject   string `json:"subject"`
	Text      string `json:"text"`
}

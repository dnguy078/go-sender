package request

import "encoding/json"

type EmailRequest struct {
	ToEmail   string `json:"toEmail"`
	FromEmail string `json:"fromEmail"`
	Subject   string `json:"subject"`
	Text      string `json:"text"`
}

func Validate(payload []byte) (EmailRequest, error) {
	e := EmailRequest{}
	if err := json.Unmarshal(payload, &e); err != nil {
		return e, err
	}

	return e, nil
}

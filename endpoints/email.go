package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Emailer struct {
}

type EmailRequest struct {
	ToEmail string `json:"email"`
}

func NewEmailer() *Emailer {
	fmt.Println("hello")
	return &Emailer{}
}

func (e *Emailer) Email(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hello world")
	req := &EmailRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println(req)

	w.WriteHeader(http.StatusOK)
}

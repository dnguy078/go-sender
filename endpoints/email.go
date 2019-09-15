package endpoints

import (
	"encoding/json"
	"net/http"
)

type EmailHandler struct {
	d dispatcher
}

type dispatcher interface {
	Dispatch() error
}

type EmailRequest struct {
	ToEmail string `json:"email"`
}

func NewEmailerHandler(dispatcher dispatcher) *EmailHandler {
	return &EmailHandler{
		d: dispatcher,
	}
}

func (e *EmailHandler) Email(w http.ResponseWriter, r *http.Request) {
	req := &EmailRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := e.d.Dispatch(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error"))
		return
	}

}

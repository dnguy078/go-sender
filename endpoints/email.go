package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/dnguy078/go-sender/request"
)

type EmailHandler struct {
	d dispatcher
}

type dispatcher interface {
	Dispatch(request.EmailRequest) error
}

func NewEmailerHandler(dispatcher dispatcher) *EmailHandler {
	return &EmailHandler{
		d: dispatcher,
	}
}

func (e *EmailHandler) Email(w http.ResponseWriter, r *http.Request) {
	req := request.EmailRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := e.d.Dispatch(req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error"))
		return
	}
}

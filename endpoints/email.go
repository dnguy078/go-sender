package endpoints

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dnguy078/go-sender/request"
)

type EmailHandler struct {
	d dispatcher
}

type dispatcher interface {
	Dispatch() error
}

func NewEmailerHandler(dispatcher dispatcher) *EmailHandler {
	return &EmailHandler{
		d: dispatcher,
	}
}

func (e *EmailHandler) Email(w http.ResponseWriter, r *http.Request) {
	req := &request.EmailRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		fmt.Println("hello", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := e.d.Dispatch(); err != nil {
		// TODO replace log with logrus and set log appropriate log level
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Service unavailable"))
		return
	}
}

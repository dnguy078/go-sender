package server

import (
	"fmt"
	"net"
	"net/http"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"

	"github.com/dnguy078/go-sender/adapter"
	"github.com/dnguy078/go-sender/config"
	"github.com/dnguy078/go-sender/endpoints"
	"github.com/dnguy078/go-sender/services"
)

type Server struct {
	router     *mux.Router
	httpServer *http.Server
	middleware *negroni.Negroni
	dispatcher services.Dispatcher

	host string
	port int
}

func (s *Server) initializeRoutes(cfg config.Config) {
	sgClient := adapter.NewSendGridClient(cfg.SendGridAPIKey)
	mgClient := adapter.NewMailgunClient()

	emailDispatcher := services.NewDispatcher("email", mgClient, sgClient)

	ee := endpoints.NewEmailerHandler(emailDispatcher)

	s.router.HandleFunc("/email", ee.Email).Methods("POST")
}

func (s *Server) initializeMiddleware() {
	// s.middleware.Use(ngLogger.NewMiddleware())

	s.middleware.UseHandler(s.router)
}

func New(cfg config.Config) *Server {
	s := &Server{
		router:     mux.NewRouter(),
		middleware: negroni.New(),
		host:       "0.0.0.0",
		port:       cfg.SenderAPIPort,
	}

	s.initializeRoutes(cfg)
	s.initializeMiddleware()

	return s
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	s.httpServer = &http.Server{Addr: addr, Handler: s.middleware}

	// start hystrix metrics server
	// separate into another function and handle error
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	go http.ListenAndServe(net.JoinHostPort("", "81"), hystrixStreamHandler)

	fmt.Printf("SenderAPI listening on %s.....\n", addr)
	if err := s.httpServer.ListenAndServe(); err != nil {
		return fmt.Errorf("error occurred when starting up supermarket %s", err)
	}

	return nil
}

// SetPrimarySender allows you to set dispatcher's primary sender; used for testing purposes
func (s *Server) SetPrimarySender(primary services.Emailer) {
	s.dispatcher.SetPrimary(primary)
}

// SetFallbackSender allows you to set dispatcher's primary sender; used for testing purposes
func (s *Server) SetFallbackSender(fallbackSender services.Emailer) {
	s.dispatcher.SetFallback(fallbackSender)
}

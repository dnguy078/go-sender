package server

import (
	"fmt"
	"net"
	"net/http"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"

	"github.com/dnguy078/go-sender/adapter"
	"github.com/dnguy078/go-sender/endpoints"
	"github.com/dnguy078/go-sender/services"
)

type Server struct {
	router     *mux.Router
	httpServer *http.Server
	middleware *negroni.Negroni

	host string
	port int
}

type ServerConfig struct {
	Port int
}

func (s *Server) initializeRoutes() {
	sgClient := adapter.NewSendGridClient("somsdfklj", "lksdfjslkdfj")
	mgClient := adapter.NewMailgunClient()

	emailDispatcher := services.NewDispatcher("email", sgClient, mgClient)

	ee := endpoints.NewEmailerHandler(emailDispatcher)

	s.router.HandleFunc("/email", ee.Email).Methods("POST")
}

func (s *Server) initializeMiddleware() {
	// s.middleware.Use(ngLogger.NewMiddleware())

	s.middleware.UseHandler(s.router)
}

func New(cfg ServerConfig) *Server {
	s := &Server{
		router:     mux.NewRouter(),
		middleware: negroni.New(),
		host:       "0.0.0.0",
		port:       cfg.Port,
	}

	s.initializeRoutes()
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
	go http.ListenAndServe(net.JoinHostPort("localhost", "81"), hystrixStreamHandler)

	fmt.Printf("SenderAPI listening on %s.....\n", addr)
	if err := s.httpServer.ListenAndServe(); err != nil {
		return fmt.Errorf("error occurred when starting up supermarket %s", err)
	}

	return nil
}

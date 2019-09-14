package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/meatballhat/negroni-logrus"
	"github.com/urfave/negroni"

	"github.com/dnguy078/go-sender/endpoints"
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
	db := marketdb.NewMarketDB()
	createProduceEndpoint := endpoints.NewCreateProduce(db)

	s.router.HandleFunc("/email", createProduceEndpoint.CreateProduce).Methods("POST")
}

// func (s *Server) initializeMiddleware() {
// 	s.middleware.Use(negronilogrus.NewMiddleware())

// 	s.middleware.UseHandler(s.router)
// }

func New(cfg ServerConfig) *Server {
	s := &Server{
		router:     mux.NewRouter(),
		middleware: negroni.New(),
		host:       "0.0.0.0",
		port:       cfg.Port,
	}

	s.initializeRoutes()
	// s.initializeMiddleware()

	return s
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	s.httpServer = &http.Server{Addr: addr, Handler: s.middleware}

	fmt.Printf("SenderAPI listening on %s.....\n", addr)
	if err := s.httpServer.ListenAndServe(); err != nil {
		return fmt.Errorf("error occurred when starting up supermarket %s", err)
	}

	return nil
}

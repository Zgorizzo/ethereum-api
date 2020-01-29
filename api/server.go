package api

import (
	"net/http"
	"time"

	"github.com/INFURA/infra-test-benjamin-mateo/config"
	"github.com/INFURA/infra-test-benjamin-mateo/node"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// Server represents the application as a struct.
// It contains the shared dependencies of the application and avoid using global state variables
type Server struct {
	router *mux.Router
	Logger *zap.SugaredLogger
	//Client instantiate client one
	client node.CustomClient
}

// NewServer bind handlers functions and set router, eth client and logger
func NewServer(logger *zap.SugaredLogger, router *mux.Router) *Server {
	s := &Server{}
	// set the logger
	s.Logger = logger
	// set the router
	s.router = router
	// enforce no cache
	s.router.Use(noCacheHeader)
	s.routes()
	return s
}

// loadClient load an ethereum client that targets the url provided.
func (s *Server) loadClient(target string) {
	s.Logger.Infof("Connecting to: %s", target)
	client, err := node.GetNewCustomClient(target)
	if err != nil {
		s.Logger.Fatal("Client error: ", err)
	}
	s.Logger.Infof("IsBidirectional  : %v", client.IsBidirectional())
	s.client = client
}

// noCacheHeader is a middleware function, to enforce no caching which will be called for each request
func noCacheHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
		next.ServeHTTP(w, r)
	})
}

// Serve the api at servingURL URL
func (s *Server) Serve(servingURL string) {
	defer s.Logger.Sync()

	// load the ethereum client
	s.loadClient(config.ReadString("NODE_URL"))

	// configure the api server
	srv := &http.Server{
		Handler: s.router,
		Addr:    servingURL,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: time.Duration(config.ReadInt("API_TIMEOUT")) * time.Second,
		ReadTimeout:  time.Duration(config.ReadInt("API_TIMEOUT")) * time.Second,
	}
	s.Logger.Info(servingURL)
	s.Logger.Fatal(srv.ListenAndServe())
}

package web

import (
	"net/http"
	"time"

	"github.com/labbsr0x/goh/gohserver"

	"github.com/abilioesteves/health-checker/checker"

	"github.com/abilioesteves/health-checker/config"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

// Server holds the information needed to run the HC Server
type Server struct {
	*config.Builder
}

// InitFromBuilder builds a Server instance
func (s *Server) InitFromBuilder(builder *config.Builder) *Server {
	s.Builder = builder

	logLevel, err := logrus.ParseLevel(s.LogLevel)
	if err != nil {
		logrus.Errorf("Not able to parse log level string. Setting default level: info.")
		logLevel = logrus.InfoLevel
	}
	logrus.SetLevel(logLevel)

	return s
}

// Run initializes the web server and its apis
func (s *Server) Run() error {
	router := mux.NewRouter().StrictSlash(true)
	router.Handle("/metrics", promhttp.Handler()).Methods("GET")
	router.Handle("/health", s.HealthHandler()).Methods("GET")

	srv := &http.Server{
		Handler:      router,
		Addr:         "0.0.0.0:" + s.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logrus.Info("Start server")
	err := srv.ListenAndServe()
	if err != nil {
		logrus.Fatal("server initialization error", err)
		return err
	}
	return nil
}

// HealthHandler defines the http handler for the self health data
func (s *Server) HealthHandler() http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		toReturn := checker.HealthCheckResponse{}

		gohserver.WriteJSONResponse(toReturn, 200, resp)
	})
}

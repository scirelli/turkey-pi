package server

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/scirelli/turkey-pi/pkg/keyboard"
	"github.com/scirelli/turkey-pi/pkg/log"
)

const (
	DEFAULT_INPUT_BUFFER_SZ uint = 500
	inputLogLength          uint = 20
)

func New(config Config, logger log.Logger, kb *keyboard.File) *Server {
	var server = Server{
		config:        config,
		logger:        logger,
		keyboardFile:  kb,
		inputBufferSz: config.InputBufferSize,
	}

	server.addr = fmt.Sprintf("%s:%d", config.Address, config.Port)
	server.registerHTTPHandlers()

	return &server
}

type Server struct {
	logger        log.Logger
	addr          string
	config        Config
	keyboardFile  *keyboard.File
	inputBufferSz uint
}

func (s *Server) Run() {
	s.logger.Infof("Listening on %s\n", s.addr)
	s.logger.Fatal(http.ListenAndServe(s.addr, nil))
}

func (s *Server) registerHTTPHandlers() {
	r := mux.NewRouter()

	s.registerStringRoutes(r.PathPrefix("/write").Subrouter())

	r.PathPrefix("/").Handler(http.FileServer(http.Dir(filepath.Join(s.config.ContentPath, "/web/static"))))

	loggedRouter := handlers.RecoveryHandler()(handlers.LoggingHandler(os.Stdout, r))
	http.Handle("/", loggedRouter)
}

/*
Notes:
	Routes are tested in the order they were added to the router. If two routes match, the first one wins:
*/
func (s *Server) registerStringRoutes(router *mux.Router) *mux.Router {
	router.Path("/string").Methods("POST").Handler(handlers.ContentTypeHandler(http.HandlerFunc(s.typeLongStringContentTypeRouterHandlerFunc), "text/plain", "application/x-www-form-urlencoded")).Name("typeLongStrings")

	return router
}

//typeLongStringContentTypeRouterHandlerFunc Route to handler based on content type. This was a quick hack to allow form submission.
func (s *Server) typeLongStringContentTypeRouterHandlerFunc(w http.ResponseWriter, r *http.Request) {
	contentType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		respondError(w, http.StatusUnsupportedMediaType, "")
		s.logger.Error(err)
		return
	}
	switch {
	case contentType == "text/plain":
		s.typeLongStringHandlerFunc(w, r)
		return
	case contentType == "application/x-www-form-urlencoded":
		s.typeLongStringFormHandlerFunc(w, r)
		return
	default: //Should never make it here since ContentTypeHandler validates the content types
		respondError(w, http.StatusUnsupportedMediaType, "")
		return
	}
}

func (s *Server) typeLongStringHandlerFunc(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var err error
	var n int
	var buf []byte = make([]byte, s.inputBufferSz)
	var totalCharRead int = 0

	for {
		if n, err = r.Body.Read(buf); err == io.EOF {
			s.logger.Debug("Reached EOF")
			break
		} else if err != nil {
			respondError(w, 503, "Failed to read input.")
			s.logger.Error(err)
			return
		}

		totalCharRead += n
		if _, err := s.keyboardFile.WriteStringDelayed(string(buf[:n])); err != nil {
			respondError(w, 502, "Failed to type message.")
			s.logger.Error(err)
			return
		}
		s.logger.Debugf("Wrote '%s'...", string(buf[:min(int(inputLogLength), n)]))
	}

	respondJSON(w, http.StatusAccepted, struct {
		Msg string `json: "msg"`
	}{
		Msg: fmt.Sprintf("Message recieved (%d char) and is being typed out", totalCharRead),
	})
}

func (s *Server) typeLongStringFormHandlerFunc(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var err error

	if err = r.ParseForm(); err != nil {
		respondError(w, 503, "Failed to read form input.")
		s.logger.Error(err)
		return
	}

	text := r.FormValue("text")
	if text == "" {
		respondError(w, http.StatusUnprocessableEntity, "Form field 'text' is required")
		s.logger.Error(err)
		return
	}
	if _, err := s.keyboardFile.WriteStringDelayed(text); err != nil {
		respondError(w, 502, "Failed to type message.")
		s.logger.Error(err)
		return
	}
	s.logger.Debugf("Form text '%s'", text)
	respondJSON(w, http.StatusAccepted, struct {
		Msg string `json: "msg"`
	}{
		Msg: fmt.Sprintf("Message recieved (%d char) and is being typed out", len(text)),
	})
}

// respondJSON makes the response with payload as json format
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

// respondError makes the error response with payload as json format
func respondError(w http.ResponseWriter, code int, message string) {
	respondJSON(w, code, map[string]string{"error": message})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

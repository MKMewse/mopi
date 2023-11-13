package server

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
)

type Server struct {
	Port      int
	Logger    *slog.Logger
	Responses map[string]Response
	Store     Storer
}

type Response struct {
	Code int         `json:"status"`
	Body interface{} `json:"body"`
	Url  string      `json:"url"`
}

func NewServer(port int, store Storer) *Server {
	return &Server{
		Port:      port,
		Logger:    slog.New(slog.NewTextHandler(os.Stdout, nil)),
		Responses: make(map[string]Response),
		Store:     store,
	}
}

func (s *Server) Run() {
	http.HandleFunc("/register", s.Register)
	http.HandleFunc("/list", s.List)
	http.HandleFunc("/reset", s.Reset)
	http.HandleFunc("/", s.MockEndpoint)
	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", s.Port), nil)
}

func (s *Server) LoadResponsesFromStore() {
	responses := s.Store.GetAll()
	for _, r := range responses {
		s.Responses[r.Url] = r
	}
}

func (s *Server) Register(w http.ResponseWriter, r *http.Request) {
	var resp Response
	err := json.NewDecoder(r.Body).Decode(&resp)

	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(400)
	}
	s.Responses[resp.Url] = resp
	go s.Store.Add(resp)
	w.WriteHeader(201)
}

func (s *Server) Reset(w http.ResponseWriter, r *http.Request) {
	s.Responses = make(map[string]Response)
	s.Store.RemoveAll()
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) List(w http.ResponseWriter, r *http.Request) {
	responsesJson, err := json.MarshalIndent(s.Responses, "", " ")
	if err != nil {
		slog.Error("Unable to marshal responses")
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(responsesJson)
}

func (s *Server) MockEndpoint(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	resp, ok := s.Responses[url]
	if ok {
		out, err := json.Marshal(resp.Body)
		if err != nil {
			slog.Error(fmt.Sprintf("Unable to unmarshall body for response %s", url))
			w.WriteHeader(500)
		} else {
			w.Write(out)
		}

	} else {
		log.Printf("No response found for url %s", url)
		w.WriteHeader(404)
	}
}

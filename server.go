// initiate comstock server instance with configuration
package main

import (
	"github.com/codegangsta/negroni"
	"net/http"
)

type Config struct {
	Port string
}

type Server struct {
	mux     *http.ServeMux
	config  Config
	negroni *negroni.Negroni
}

func NewServer(config Config) *Server {

	mux := http.NewServeMux()
	mux.HandleFunc("/loginAs", LoginAsHandler)
	mux.HandleFunc("/checkSession", CheckSessionHandler)
	mux.HandleFunc("/list", ListHandler)
	mux.HandleFunc("/registerUser", RegistUserHandler)
	mux.HandleFunc("/postOne", PostHistoryHandler)
	mux.HandleFunc("/fetchOne", FetchHandler)
	mux.HandleFunc("/removeOne", RemoveOneHandler)
	//	mux.HandleFunc("/postChunkCommands", PostChunkCommandsHandler)
	//	mux.HandleFunc("/trans", TransHandler)
	mux.HandleFunc("/js/angular.min.js.map", func(w http.ResponseWriter, req *http.Request) {
		return
	})

	return &Server{
		negroni: negroni.Classic(),
		mux:     mux,
		config:  config,
	}
}

func (s *Server) Start() {

	s.negroni.UseHandler(s.mux)
	s.negroni.Run(":" + s.config.Port)
}

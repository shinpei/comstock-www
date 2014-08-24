package main

import (
	"github.com/codegangsta/negroni"
	"log"
	"net/http"
	"os"
)

const (
	Version string = "0.1.4-devel"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/loginAs", LoginAsHandler)
	mux.HandleFunc("/checkSession", CheckSessionHandler)
	mux.HandleFunc("/list", ListHandler)
	mux.HandleFunc("/registerUser", RegistUserHandler)
	mux.HandleFunc("/postCommand", PostCommandHandler)
	mux.HandleFunc("/fetchCommandFromNumber", FetchHandler)
	mux.HandleFunc("/removeOne", RemoveOneHandler)

	n := negroni.Classic()
	n.UseHandler(mux)
	port := ""
	if os.Getenv("PORT") == "" {
		port = "5000"
	} else {
		port = os.Getenv("PORT")
	}
	log.Printf("type of nil:%T\n", nil)
	log.Println("Start webserver")
	n.Run(":" + port)
}

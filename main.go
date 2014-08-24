package main

import (
	"github.com/codegangsta/negroni"
	cmodel "github.com/shinpei/comstock/model"
	"log"
	"net/http"
	"net/url"
	"os"
)

const (
	Version string = "0.1.4-devel"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/loginAs", LoginAsHandler)

	mux.HandleFunc("/checkSession", func(w http.ResponseWriter, req *http.Request) {
		session, db := getSessionAndDB()
		defer session.Close()
		// make sure param exists
		m, _ := url.ParseQuery(req.URL.RawQuery)
		if m["authinfo"] == nil {
			// error
			log.Println("Error, check session requires param")
			http.Error(w, "session check needs parameters", http.StatusBadRequest)
			return
		}
		_, err := GetUserSession(db, m["authinfo"][0])
		if err == cmodel.ErrSessionNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
		}
		// if pass reaches here, session found. do nothing.
	})

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
	n.Run(":" + port)
}

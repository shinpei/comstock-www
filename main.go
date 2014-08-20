package main

import (
	//	"encoding/json"
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/shinpei/comstock-www/model"
	cmodel "github.com/shinpei/comstock/model"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/loginAs", func(w http.ResponseWriter, req *http.Request) {
		session, db := getSessionAndDB()

		// make sure param exists
		params, _ := url.ParseQuery(req.URL.RawQuery)
		if params["mail"] == nil || params["password"] == nil {
			// error
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		fmt.Printf("login request mail:%#v, %#v", params["mail"], params["mail"][0])
		s, err := LoginAs(db, model.CreateLoginRequest(params["mail"][0], params["password"][0]))
		if err == cmodel.ErrUserNotFound || err == cmodel.ErrIncorrectPassword {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		/*
			resJson, err := json.Marshal(s.Token)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		*/
		w.Header().Set("Content-type", "application/json")
		if err == cmodel.ErrAlreadyLogin {
			w.WriteHeader(http.StatusConflict)

		}

		w.Write([]byte(s.Token))
		session.Close()
	})

	mux.HandleFunc("/checkSession", func(w http.ResponseWriter, req *http.Request) {
		session, db := getSessionAndDB()

		// make sure param exists
		m, _ := url.ParseQuery(req.URL.RawQuery)
		if m["authinfo"] != nil {
			err := CheckSession(db, m["authinfo"][0])
			if err == cmodel.ErrSessionNotFound {
				http.Error(w, err.Error(), http.StatusNotFound)
			}
			// if pass reaches here, session found. do nothing.
		} else {
			// error
			fmt.Println("Error, check session requires param")
			http.Error(w, "session check needs parameters", http.StatusBadRequest)
			return
		}
		session.Close()
	})

	mux.HandleFunc("/list", func(w http.ResponseWriter, req *http.Request) {
		session, db := getSessionAndDB()
		m, _ := url.ParseQuery(req.URL.RawQuery)
		if m["authinfo"] != nil {
			err := ListCommands(db, m["authinfo"][0])
			if err != cmodel.ErrSessionNotFound {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
		} else {
			// error
			log.Println("Error, check session requires param")
			http.Error(w, "Session check needs parameters", http.StatusBadGateway)
			return
		}
		session.Close()
	})

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

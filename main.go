package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/shinpei/comstock-www/model"
	cmodel "github.com/shinpei/comstock/model"
	"net/http"
	"net/url"
	"os"
)

func main() {
	mux := http.NewServeMux()
	_, db := getSessionAndDB()

	mux.HandleFunc("/loginAs", func(w http.ResponseWriter, req *http.Request) {
		// make sure param exists
		params, _ := url.ParseQuery(req.URL.RawQuery)
		if params["mail"] == nil || params["password"] == nil {
			// error
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		res, err := LoginAs(db, model.CreateLoginRequest(params["mail"][0], params["password"][0]))
		if err == cmodel.ErrUserNotFound || err == cmodel.ErrIncorrectPassword {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		resJson, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-type", "application/json")
		w.Write(resJson)
	})

	mux.HandleFunc("/checkSession", func(w http.ResponseWriter, req *http.Request) {
		// make sure param exists
		m, _ := url.ParseQuery(req.URL.RawQuery)
		if m["authinfo"] != nil {
			err := CheckSession(db, m["authinfo"][0])
			if err != cmodel.ErrSessionNotFound {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
		} else {
			// error
			fmt.Println("Error, check session requires query")
			http.Error(w, "session check needs parameters", http.StatusBadRequest)
			return
		}
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

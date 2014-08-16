package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/shinpei/comstock-www/model"
	"net/http"
	"net/url"
	"os"
)

func main() {
	mux := http.NewServeMux()
	_, db := getSessionAndDB()

	mux.HandleFunc("/loginAs", func(w http.ResponseWriter, req *http.Request) {
		res, err := LoginAs(db, model.CreateLoginInfo("hoge", "hi"))
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
			CheckSession(db, m["authinfo"][0])
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

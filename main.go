package main

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"net/http"
	"net/url"
	"os"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/loginAs", func(w http.ResponseWriter, req *http.Request) {
		u, err := url.Parse(os.Getenv("MONGOHQ_URL"))
		if err != nil {
			panic("error")
		}
		fmt.Fprintf(w, "scheme:%s, user:%s, username:%s ", u.Scheme, u.User, u.User.Username())
	})

	n := negroni.Classic()
	n.UseHandler(mux)
	port := ""
	if os.Getenv("PORT") == "" {
		port = "5000"
	} else {
		port = os.Getenv("PORT")
	}
	//session, db := getSessionAndDB()

	n.Run(":" + port)
}

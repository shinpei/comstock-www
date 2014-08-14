package main

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/loginAs", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "hi")

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

package main

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/shinpei/comstock/model"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/loginAs", func(w http.ResponseWriter, req *http.Request) {
		session, db := getSessionAndDB()
		fmt.Fprintf(w, "scheme:%#v, %#v<p>", session, db)
		c := db.C("command")
		iter := c.Find(nil).Iter()
		var result model.Command
		for iter.Next(&result) {
			fmt.Fprintf(w, "%#v<br>", result.Cmd)
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
	//session, db := getSessionAndDB()

	n.Run(":" + port)
}

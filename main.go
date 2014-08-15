package main

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/shinpei/comstock/model"
	"log"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()
	_, db := getSessionAndDB()

	mux.HandleFunc("/loginAs", func(w http.ResponseWriter, req *http.Request) {
		c := db.C("command")
		iter := c.Find(nil).Iter()
		log.Println(iter)
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
	n.Run(":" + port)
}

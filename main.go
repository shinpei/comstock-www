package main

import (
	"github.com/codegangsta/negroni"
	"github.com/shinpei/comstock-www/model"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()
	_, db := getSessionAndDB()

	mux.HandleFunc("/loginAs", func(w http.ResponseWriter, req *http.Request) {
		LoginAs(db, model.CreateLoginInfo("hoge", "hi"))

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

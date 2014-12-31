/* Comstock server
 *
 */
package main

import (
	//	"github.com/pkg/profile"
	"os"
)

const (
	Version string = "0.2.0-pre2"
)

func main() {
	port := ""
	//	defer profile.Start().Stop()
	if os.Getenv("PORT") == "" {
		port = "5000"
	} else {
		port = os.Getenv("PORT")
	}

	server := NewServer(Config{
		Port: port,
	})

	server.Start()
}

/* Comstock server
 *
 */
package main

import (
	"os"
)

const (
	Version string = "0.2.0-pre"
)

func main() {
	port := ""
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

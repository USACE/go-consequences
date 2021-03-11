package main

import (
	"log"
	"net/http"

	"github.com/USACE/go-consequences/compute"
)

//Config describes the configuration settings for go-consequences.
type Config struct {
	SkipJWT       bool
	LambdaContext bool
	DBUser        string
	DBPass        string
	DBName        string
	DBHost        string
	DBSSLMode     string
}

func main() {
	var cfg Config
	if cfg.LambdaContext {
		log.Print("starting server; Running On AWS LAMBDA")
		//lambda.Start(HandleRequestArgs)
	} else {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			params := r.URL.Query()
			fp, fpPresent := params["FilePath"]
			if !fpPresent {
				http.Error(w, "No FilePath argument", http.StatusNotFound)
			} else {
				if len(fp[0]) == 0 {
					//should have better error checking...
					http.Error(w, "Invalid FilePath argument", http.StatusNotFound)
				} else {
					//s, _ := compute.FromFile(fp[0])
					//fmt.Fprintf(w, s)
					compute.StreamFromFile(fp[0], w)
				}
			}
		})
		log.Print("Not on Lambda")
		log.Print("starting local server")
		log.Fatal(http.ListenAndServe("localhost:3030", nil))
	}
}

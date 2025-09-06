package main

import (
	"log"
	"net/http"

	"github.com/AJMerr/MSE/pkg/store"
	"github.com/AJMerr/NYMToDo/pkg/api"
	"github.com/AJMerr/gonk/pkg/router"
	"github.com/AJMerr/parsec/pkg/parsec"
)

func main() {
	s := store.NewStore()
	h := api.New(s)

	r := router.NewRouter()
	h.RegisterRouter(r)

	addr := ":8080"

	go func() {
		log.Printf("NYMToDo API listening on %s", addr)
		if err := http.ListenAndServe(addr, r); err != nil {
			log.Fatal(err)
		}
	}()

	ph, err := parsec.HandlerFromFile("./parsec.json")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("UI proxy listening on :5173")
	log.Fatal(http.ListenAndServe(":5173", ph))
}

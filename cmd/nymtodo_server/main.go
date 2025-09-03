package main

import (
	"log"
	"net/http"

	"github.com/AJMerr/MSE/pkg/store"
	"github.com/AJMerr/NYMToDo/pkg/api"
	"github.com/AJMerr/gonk/pkg/router"
)

func main() {
	s := store.NewStore()
	h := api.New(s)

	r := router.NewRouter()
	h.RegisterRouter(r)

	addr := ":8080"
	log.Printf("NYMToDo listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}

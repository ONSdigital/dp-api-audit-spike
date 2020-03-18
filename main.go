package main

import (
	"fmt"
	"net/http"

	"github.com/ONSdigital/dp-api-audit-spike/auditing"
	"github.com/ONSdigital/dp-api-audit-spike/handlers"
	"github.com/gorilla/mux"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err.Error())
	}
}

func run() error {
	r := mux.NewRouter()

	auditor := &auditing.Stub{}

	foo := handlers.Foo(auditor)
	r.Handle("/foo", foo)

	bar := handlers.Bar(auditor)
	r.Handle("/bar/{name}", bar)

	return http.ListenAndServe(":8080", r)
}

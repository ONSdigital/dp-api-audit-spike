package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ONSdigital/dp-api-audit-spike/handlers"
	"github.com/ONSdigital/go-ns/common"
	"github.com/fatih/color"
	"github.com/gorilla/mux"
)

var yellowOut = color.New(color.FgHiYellow)

type AuditorStub struct {
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err.Error())
	}
}

func run() error {
	r := mux.NewRouter()

	auditor := &AuditorStub{}

	foo := handlers.Foo(auditor)
	r.Handle("/foo", foo)

	bar := handlers.Bar(auditor)
	r.Handle("/bar/{name}", bar)

	return http.ListenAndServe(":8080", r)
}

func (s *AuditorStub) Record(ctx context.Context, action string, result string, params common.Params) error {
	yellowOut.Printf("auditing service: recording event  - action: %s, status: %s, params %+v\n", action, result, params)
	return nil
}

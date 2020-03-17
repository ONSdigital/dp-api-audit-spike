package main

import (
	"fmt"
	"net/http"

	"github.com/ONSdigital/dp-api-audit-spike/auditing"
	"github.com/ONSdigital/go-ns/audit"
	"github.com/ONSdigital/go-ns/common"
	"github.com/gorilla/mux"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err.Error())
	}
}

func run() error {
	auditor := &audit.NopAuditor{}

	r := mux.NewRouter()
	r.Handle("/foo", fooHandler(auditor))
	r.Handle("/bar", barHandler(auditor))

	return http.ListenAndServe(":8080", r)
}

func fooHandler(auditor audit.AuditorService) http.Handler {
	return &auditing.Handler{
		Action:  "GET Foo",
		Auditor: auditor,
		GetParams: func(r *http.Request) common.Params {
			return nil
		},
		Successful: func(w *auditing.ResponseWriter) bool {
			return w.Status == 200
		},
		Handler: func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("inside Foo handler")
			w.Write([]byte("Foo"))
		},
	}
}

func barHandler(auditor audit.AuditorService) http.Handler {
	return &auditing.Handler{
		Action:  "GET Bar",
		Auditor: auditor,
		GetParams: func(r *http.Request) common.Params {
			return nil
		},
		Successful: func(w *auditing.ResponseWriter) bool {
			return w.Status == 200
		},
		Handler: func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("inside Bar handler")
			w.Write([]byte("Bar"))
		},
	}
}

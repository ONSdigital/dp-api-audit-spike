package handlers

import (
	"net/http"

	"github.com/ONSdigital/dp-api-audit-spike/auditing"
	"github.com/ONSdigital/go-ns/audit"
	"github.com/ONSdigital/go-ns/common"
	"github.com/fatih/color"
	"github.com/gorilla/mux"
)

var greenOut = color.New(color.FgHiGreen)

func Foo(auditor audit.AuditorService) http.Handler {
	fooAuditParams := func(r *http.Request) common.Params {
		return nil
	}

	return auditing.Wrap(fooHandleFunc(), "GET foo", auditor, fooAuditParams, 200)
}

func Bar(auditor audit.AuditorService) http.Handler {
	barAuditParams := func(r *http.Request) common.Params {
		vars := mux.Vars(r)
		return common.Params{"name": vars["name"]}
	}

	return auditing.Wrap(barHandleFunc(), "GET bar", auditor, barAuditParams, 200)
}

func fooHandleFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		greenOut.Println("foo handler: handling request for get foo")
		w.Write([]byte("Foo"))
	}
}

func barHandleFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		greenOut.Println("bar handler: handling request for get bar")
		w.Write([]byte("Bar"))
	}
}

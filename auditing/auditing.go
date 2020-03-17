package auditing

import (
	"fmt"
	"net/http"

	"github.com/ONSdigital/go-ns/audit"
	"github.com/ONSdigital/go-ns/common"
)

type Handler struct {
	// The action being carried out.
	Action string
	// The auditing service
	Auditor audit.AuditorService
	// Get the audit params from the request
	GetParams func(r *http.Request) common.Params
	// Was the action successful
	Successful func(w *ResponseWriter) bool
	// The handler for this request
	Handler http.HandlerFunc
}

type ResponseWriter struct {
	http.ResponseWriter
	Status int
}

func (s *ResponseWriter) WriteHeader(statusCode int) {
	s.Status = statusCode
	s.ResponseWriter.WriteHeader(statusCode)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("inside audit wrapper")
	ctx := r.Context()
	// get the audit parameters.
	params := h.GetParams(r)

	// audit the action is being attempted
	if err := h.Auditor.Record(ctx, h.Action, audit.Attempted, params); err != nil {
		fmt.Println("audit action attempted failed")
		http.Error(w, "audit action attempted failed", 500)
		return
	}

	respW := &ResponseWriter{
		Status:         0,
		ResponseWriter: w,
	}

	// call the http handler for this request
	fmt.Println("calling wrapped handler")
	h.Handler.ServeHTTP(respW, r)

	// check the result handling the request
	result := audit.Unsuccessful
	if h.Successful(respW) {
		result = audit.Successful
	}

	// audit the outcome of the action
	fmt.Println("audit action result")
	err := h.Auditor.Record(ctx, h.Action, result, params)
	if err != nil {
		http.Error(w, "audit action result failed", 500)
	}
}

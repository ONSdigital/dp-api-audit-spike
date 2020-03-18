package auditing

import (
	"context"
	"net/http"

	"github.com/ONSdigital/go-ns/common"
	"github.com/fatih/color"
)

var (
	blueOut = color.New(color.FgHiBlue)
	redOut  = color.New(color.FgHiRed)
)

type Service interface {
	Record(ctx context.Context, action string, result string, params common.Params) error
}

type Handler struct {
	// The HTTP handler for this request
	Handler http.HandlerFunc
	// The action being carried out.
	Action string
	// The auditing service
	Auditor Service
	// Get the audit params from the request
	GetAuditParams GetAuditParamsFunc
	// The HTTP status code if the action was successful
	SuccessStatus int
}

type GetAuditParamsFunc func(r *http.Request) common.Params

type ResponseWriter struct {
	http.ResponseWriter
	Status int
}

// Wrap create a http Handler that wraps around the provided handler capturing auditing data.
func Wrap(h http.HandlerFunc, action string, auditor Service, getParamsFunc GetAuditParamsFunc, successStatus int) http.Handler {
	return &Handler{
		Handler:        h,
		Action:         action,
		Auditor:        auditor,
		SuccessStatus:  successStatus,
		GetAuditParams: getParamsFunc,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	blueOut.Println("audit handler: recording action attempted")
	ctx := r.Context()
	// get the audit parameters.
	params := h.GetAuditParams(r)

	// audit the action is being attempted
	if err := h.Auditor.Record(ctx, h.Action, "attempted", params); err != nil {
		redOut.Println("audit action attempted failed")
		http.Error(w, "audit action attempted failed", 500)
		return
	}

	respW := &ResponseWriter{
		Status:         200,
		ResponseWriter: w,
	}

	// call the wrapped http handler to execute the requested action
	blueOut.Println("audit handler: invoking wrapped handler")
	h.Handler.ServeHTTP(respW, r)

	// check the result of handling the request
	result := "unsuccessful"
	if h.SuccessStatus == respW.Status {
		result = "successful"
	}

	// audit the outcome of the action
	blueOut.Println("audit handler: recording action result")
	err := h.Auditor.Record(ctx, h.Action, result, params)
	if err != nil {
		redOut.Println("audit action result failed")
		http.Error(w, "audit action result failed", 500)
	}
	blueOut.Println("audit handler: request complete\n")
}

func (s *ResponseWriter) WriteHeader(statusCode int) {
	s.Status = statusCode
	s.ResponseWriter.WriteHeader(statusCode)
}

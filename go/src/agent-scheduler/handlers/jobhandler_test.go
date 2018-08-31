package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateAndSubmitJobScript(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		code int
	}{
		{"base-case", args{httptest.NewRecorder(), httptest.NewRequest("POST", "/api/v1/job/request/submit", strings.NewReader(`{"JobId":"1"}`))}, http.StatusOK},
		{"unprocesseable", args{httptest.NewRecorder(), httptest.NewRequest("POST", "/api/v1/job/request/submit", strings.NewReader(`{abc}{fgh}{143}`))}, http.StatusUnprocessableEntity},
		{"bad-method", args{httptest.NewRecorder(), httptest.NewRequest("GET", "/api/v1/job/request/submit", nil)}, http.StatusMethodNotAllowed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CreateAndSubmitJobScript(tt.args.w, tt.args.r)
			w := tt.args.w.(*httptest.ResponseRecorder)
			if got, want := w.Code, tt.code; got != want {
				t.Errorf("got %d; want %d", got, want)
			}
		})
	}
}

func TestTerminateJobHandler(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		code int
	}{
		{"base-case", args{httptest.NewRecorder(), httptest.NewRequest("POST", "/api/v1/job/request/terminate", strings.NewReader(`{"JobId":"1"}`))}, http.StatusOK},
		{"unprocesseable", args{httptest.NewRecorder(), httptest.NewRequest("POST", "/api/v1/job/request/terminate", strings.NewReader(`{abc}{fgh}{143}`))}, http.StatusUnprocessableEntity},
		{"case-no-jobid", args{httptest.NewRecorder(), httptest.NewRequest("POST", "/api/v1/job/request/terminate", strings.NewReader(`{"Job": "M"}`))}, http.StatusBadRequest},
		{"bad-method", args{httptest.NewRecorder(), httptest.NewRequest("GET", "/api/v1/job/request/terminate/jobid/1", nil)}, http.StatusMethodNotAllowed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			TerminateJobHandler(tt.args.w, tt.args.r)
			w := tt.args.w.(*httptest.ResponseRecorder)
			if got, want := w.Code, tt.code; got != want {
				t.Errorf("got %d; want %d", got, want)
			}
		})
	}
}

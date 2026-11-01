package prtg

import (
	"context"
	"crypto/x509"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRejectsUntrustedPRTGTLSCertificate(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	api := NewApi(server.URL, "token", time.Minute, time.Minute)
	_, err := api.baseExecuteRequest("getstatus.htm", nil)
	if err == nil {
		t.Fatal("expected an untrusted TLS certificate to be rejected")
	}
	var unknownAuthority x509.UnknownAuthorityError
	if !errors.As(err, &unknownAuthority) {
		t.Fatalf("expected an unknown certificate authority error, got %v", err)
	}
}

func TestGetStatusListPropagatesContextCancellation(t *testing.T) {
	requestStarted := make(chan struct{})
	requestCanceled := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(requestStarted)
		<-r.Context().Done()
		close(requestCanceled)
	}))
	defer server.Close()

	api := NewApi(server.URL, "token", time.Minute, time.Minute)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	result := make(chan error, 1)
	go func() {
		_, err := api.GetStatusList(ctx)
		result <- err
	}()

	select {
	case <-requestStarted:
	case <-time.After(time.Second):
		t.Fatal("HTTP request did not start")
	}
	cancel()

	select {
	case err := <-result:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("GetStatusList() error = %v, want context canceled", err)
		}
	case <-time.After(time.Second):
		t.Fatal("GetStatusList did not return after context cancellation")
	}

	select {
	case <-requestCanceled:
	case <-time.After(time.Second):
		t.Fatal("HTTP request context was not canceled")
	}
}

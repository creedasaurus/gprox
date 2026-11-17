package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestReverseProxyForwardsRequestAndResponse(t *testing.T) {
	proxy := newReverseProxy("backend.example", 8123)
	called := false
	proxy.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		called = true
		if got, want := r.URL.String(), "http://backend.example:8123/items/a%2Fb?filter=one&filter=two"; got != want {
			t.Errorf("target URL = %q, want %q", got, want)
		}
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", got)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		if got, want := string(body), `{"name":"example"}`; got != want {
			t.Errorf("request body = %q, want %q", got, want)
		}
		return &http.Response{
			StatusCode: http.StatusCreated,
			Header:     http.Header{"Location": {"/items/42"}},
			Body:       io.NopCloser(strings.NewReader("created")),
		}, nil
	})

	req := httptest.NewRequest(http.MethodPost, "https://localhost:9001/items/a%2Fb?filter=one&filter=two", strings.NewReader(`{"name":"example"}`))
	req.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	proxy.ServeHTTP(response, req)

	if !called {
		t.Fatal("upstream transport was not called")
	}
	if response.Code != http.StatusCreated {
		t.Errorf("response status = %d, want 201", response.Code)
	}
	if got := response.Header().Get("Location"); got != "/items/42" {
		t.Errorf("Location = %q, want /items/42", got)
	}
	if got := response.Body.String(); got != "created" {
		t.Errorf("response body = %q, want created", got)
	}
}

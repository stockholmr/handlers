package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDefaultCORSHandlerReturnsOk(t *testing.T) {
	methods := []string{"GET", "HEAD", "POST"}

	for _, method := range methods {
		r := newRequest(method, "http://www.example.com/")
		rr := httptest.NewRecorder()

		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

		CORS()(testHandler).ServeHTTP(rr, r)

		if status := rr.Code; status != http.StatusOK {
			t.Fatalf("bad status: got %v want %v for method %s", status, http.StatusFound, method)
		}
	}
}

func TestCORSHandlerIgnoreOptionsFallsThrough(t *testing.T) {
	r := newRequest("OPTIONS", "http://www.example.com/")
	r.Header.Set("Origin", r.URL.String())

	rr := httptest.NewRecorder()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})

	CORS(IgnoreOptions())(testHandler).ServeHTTP(rr, r)

	if status := rr.Code; status != http.StatusTeapot {
		t.Fatalf("bad status: got %v want %v", status, http.StatusTeapot)
	}
}

func TestCORSHandlerSetsExposedHeaders(t *testing.T) {
	methods := []string{"GET", "HEAD", "POST"}

	for _, method := range methods {
		// Test default configuration.
		r := newRequest(method, "http://www.example.com/")
		r.Header.Set("Origin", r.URL.String())

		rr := httptest.NewRecorder()

		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

		CORS(ExposedHeaders([]string{"X-CORS-TEST"}))(testHandler).ServeHTTP(rr, r)

		if status := rr.Code; status != http.StatusOK {
			t.Fatalf("bad status: got %v want %v for method %s", status, http.StatusOK, method)
		}

		header := rr.HeaderMap.Get(corsExposeHeadersHeader)
		if header != "X-Cors-Test" {
			t.Fatalf("bad header: expected X-Cors-Test header, got empty header for method %s.", method)
		}
	}
}

func TestCORSHandlerUnsetRequethMethodForPreflightBadRequest(t *testing.T) {
	r := newRequest("OPTIONS", "http://www.example.com/")
	r.Header.Set("Origin", r.URL.String())

	rr := httptest.NewRecorder()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	CORS(AllowedMethods([]string{"DELETE"}))(testHandler).ServeHTTP(rr, r)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Fatalf("bad status: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestCORSHandlerAllowedMethodForPreflight(t *testing.T) {
	r := newRequest("OPTIONS", "http://www.example.com/")
	r.Header.Set("Origin", r.URL.String())
	r.Header.Set(corsRequestMethodHeader, "DELETE")

	rr := httptest.NewRecorder()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	CORS(AllowedMethods([]string{"DELETE"}))(testHandler).ServeHTTP(rr, r)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("bad status: got %v want %v", status, http.StatusOK)
	}

	header := rr.HeaderMap.Get(corsAllowMethodsHeader)
	if header != "DELETE" {
		t.Fatalf("bad header: expected DELETE method header, got empty header.")
	}
}

func TestCORSHandlerAllowedHeaderForPreflight(t *testing.T) {
	r := newRequest("OPTIONS", "http://www.example.com/")
	r.Header.Set("Origin", r.URL.String())
	r.Header.Set(corsRequestMethodHeader, "POST")
	r.Header.Set(corsRequestHeadersHeader, "Content-Type")

	rr := httptest.NewRecorder()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	CORS(AllowedHeaders([]string{"Content-Type"}))(testHandler).ServeHTTP(rr, r)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("bad status: got %v want %v", status, http.StatusOK)
	}

	header := rr.HeaderMap.Get(corsAllowHeadersHeader)
	if header != "Content-Type" {
		t.Fatalf("bad header: expected Content-Type header, got empty header.")
	}
}

func TestCORSHandlerMaxAgeForPreflight(t *testing.T) {
	r := newRequest("OPTIONS", "http://www.example.com/")
	r.Header.Set("Origin", r.URL.String())
	r.Header.Set(corsRequestMethodHeader, "POST")

	rr := httptest.NewRecorder()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	CORS(MaxAge(3500))(testHandler).ServeHTTP(rr, r)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("bad status: got %v want %v", status, http.StatusOK)
	}

	header := rr.HeaderMap.Get(corsMaxAgeHeader)
	if header != "600" {
		t.Fatalf("bad header: expected %s to be %s, got %s.", corsMaxAgeHeader, "600", header)
	}
}

func TestCORSHandlerAllowedCredentials(t *testing.T) {
	r := newRequest("GET", "http://www.example.com/")
	r.Header.Set("Origin", r.URL.String())

	rr := httptest.NewRecorder()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	CORS(AllowCredentials())(testHandler).ServeHTTP(rr, r)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("bad status: got %v want %v", status, http.StatusOK)
	}

	header := rr.HeaderMap.Get(corsAllowCredentialsHeader)
	if header != "true" {
		t.Fatalf("bad header: expected %s to be %s, got %s.", corsAllowCredentialsHeader, "true", header)
	}
}

func TestCORSHandlerMultipleAllowOriginsSetsVaryHeader(t *testing.T) {
	r := newRequest("GET", "http://www.example.com/")
	r.Header.Set("Origin", r.URL.String())

	rr := httptest.NewRecorder()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	CORS(AllowedOrigins([]string{r.URL.String(), "http://google.com"}))(testHandler).ServeHTTP(rr, r)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("bad status: got %v want %v", status, http.StatusOK)
	}

	header := rr.HeaderMap.Get(corsVaryHeader)
	if header != corsOriginHeader {
		t.Fatalf("bad header: expected %s to be %s, got %s.", corsVaryHeader, corsOriginHeader, header)
	}
}

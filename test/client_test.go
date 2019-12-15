package test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alandiegosantos/http-random-stress/internal"
	"gotest.tools/assert"
)

func NewHttpServer() *httptest.Server {

	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		if req.URL.String() == "/" {
			// Send response to be tested
			rw.Write([]byte("OK"))
		} else if req.URL.String() == "/error" {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("Error"))
		} else {
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte("Not Found"))
		}

	}))

}

func TestClientHttp2(t *testing.T) {

	// Start a local HTTP server
	server := NewHttpServer()
	// Close the server when test finishes
	defer server.Close()

	// Testing http 200 counter

	client := internal.NewHTTPClient()

	httpRequestsCount := internal.GetHttpReqCounter(2, server.URL)

	internal.DoRequest(context.Background(), client, server.URL, "GET")

	assert.Equal(t, internal.GetHttpReqCounter(2, server.URL), httpRequestsCount+1)

	// Testing http 500 counter

	httpRequestsCount = internal.GetHttpReqCounter(500, server.URL+"/error")

	internal.DoRequest(context.Background(), client, server.URL+"/error", "GET")

	assert.Equal(t, internal.GetHttpReqCounter(5, server.URL+"/error"), httpRequestsCount+1)

	// Testing http 404 counter

	httpRequestsCount = internal.GetHttpReqCounter(404, server.URL+"/notfound")

	internal.DoRequest(context.Background(), client, server.URL+"/notfound", "GET")

	assert.Equal(t, internal.GetHttpReqCounter(4, server.URL+"/notfound"), httpRequestsCount+1)

}

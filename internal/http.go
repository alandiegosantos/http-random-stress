package internal

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"net"
	"net/http"
	"time"
)

func NewHTTPClient() *http.Client {

	dialer := &net.Dialer{ // use DialContext here
		Timeout:   10 * time.Second,
		KeepAlive: 1 * time.Minute,
	}

	transport := &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		DialContext:         dialer.DialContext,
		TLSHandshakeTimeout: 10 * time.Second,
		MaxIdleConns:        3000,
		MaxIdleConnsPerHost: 3000,
		IdleConnTimeout:     60 * time.Second,
	}

	return &http.Client{
		Transport: transport,
	}

}

func DoRequest(ctx context.Context, client *http.Client, url string, method string) {

	req, err := http.NewRequestWithContext(ctx, method, url, nil)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	resp, err := client.Do(req)

	if err == nil {
		resp.Body.Close()

		httpStatusCodeClass := int(resp.StatusCode / 100)

		IncrementHttpReqCounter(httpStatusCodeClass, url)

	} else {
		fmt.Printf("Error: %v\n", err)
	}

}

func StartWorker(ctx context.Context, rate float64, url string) map[int]int {

	client := NewHTTPClient()

	for {

		select {
		case <-ctx.Done():
			return nil
		default:

			randNum := -math.Log(1.0-rand.Float64()) / rate

			time.Sleep(time.Duration(randNum) * time.Second)

			go DoRequest(ctx, client, url, "GET")

		}

	}

}

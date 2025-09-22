package http_client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

var gTr = &http.Transport{
	ResponseHeaderTimeout: time.Second * 3,
	Proxy:                 http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}).DialContext,
	ForceAttemptHTTP2:     true,
	MaxIdleConns:          200,
	MaxIdleConnsPerHost:   100,
	IdleConnTimeout:       30 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}
var gClient = &http.Client{
	Transport: gTr,
	Timeout:   3 * time.Second,
}

func HttpDo[T any](method string, header map[string]string, url string, params map[string]any, res *T) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	paramBytes, _ := json.Marshal(params)
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(paramBytes))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range header {
		req.Header.Set(k, v)
	}

	resp, err := gClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	log.Printf("http response: status=%s url=%s param=%s body=%s", resp.Status, url, string(paramBytes), string(body))
	if res != nil {
		return json.Unmarshal(body, &res)
	}
	return
}

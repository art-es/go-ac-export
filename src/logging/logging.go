package logging

import (
	"fmt"
	"net/http"
)

func RequestInfo(resp *http.Response, content []byte) {
	fmt.Printf(
		"\n\nSend request\nMethod: %s\nURI: %s\nPayload: %#v\nStatusCode: %d\nStatusMessage: %s\nBody: %s\n",
		resp.Request.Method,
		resp.Request.URL.String(),
		resp.Request.Body,
		resp.StatusCode,
		http.StatusText(resp.StatusCode),
		string(content),
	)
}

func RequestInfoWithoutBody(resp *http.Response) {
	fmt.Printf(
		"\n\nSend request\nMethod: %s\nURI: %s\nHeaders: %#v\nPayload: %#v\nStatusCode: %d\nStatusMessage: %s\n",
		resp.Request.Method,
		resp.Request.URL.String(),
		resp.Request.Header,
		resp.Request.Body,
		resp.StatusCode,
		http.StatusText(resp.StatusCode),
	)
}

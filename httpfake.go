package go_httpfake

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
)

// RequestResponse is the structure where routes are saving
type RequestResponse struct {
	returnStatus  int
	returnBody    string
	returnHeaders map[string]string

	requestParams map[string]string
}

// HTTPFake is the workhorse used to add, delete routes and to start a fake server
type HTTPFake struct {
	started bool
	server  *httptest.Server
	routes  map[string]map[string]RequestResponse
}

// New return an instance of HTTPFake
func New() *HTTPFake {
	return &HTTPFake{started: false}
}

// AddRoute is using to add a new route
func (h *HTTPFake) AddRoute(route string, method string, requestParams map[string]string, returnStatus int, returnBody string, returnHeaders map[string]string) bool {
	if h.started {
		return false
	}

	if len(h.routes) == 0 {
		h.routes = map[string]map[string]RequestResponse{}
	}

	if _, ok := h.routes[route]; !ok {
		h.routes[route] = map[string]RequestResponse{}
	}

	h.routes[route][method] = RequestResponse{
		returnStatus:  returnStatus,
		returnBody:    returnBody,
		returnHeaders: returnHeaders,

		requestParams: requestParams,
	}

	return true
}

// DelRoute is using to delete a route
func (h *HTTPFake) DelRoute(route string, method string) bool {
	if h.started {
		return false
	}

	if _, ok := h.routes[route]; !ok {
		return false
	}

	if _, ok := h.routes[route][method]; !ok {
		return false
	}

	delete(h.routes[route], method)

	if len(h.routes[route]) == 0 {
		delete(h.routes, route)
	}

	return true
}

// Start is using to start a new fake server
func (h *HTTPFake) Start() *httptest.Server {
	h.started = true

	h.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := strings.Replace(r.URL.String(), "/", "", 1)
		if _, ok := h.routes[route]; !ok {
			h.notFound(w, map[string]string{})

			return
		}

		if _, ok := h.routes[route][r.Method]; !ok {
			h.notFound(w, map[string]string{})

			return
		}

		requestResponse := h.routes[route][r.Method]

		// Convert body parameters to map
		contents, err := ioutil.ReadAll(r.Body)
		if err != nil {
			h.setResponse(w, 500, fmt.Sprintf(`{"error": "%s"}`, err.Error()), requestResponse.returnHeaders)

			return
		}

		var jsonParams map[string]interface{}
		_ = json.Unmarshal(contents, &jsonParams)

		// Check if expected parameters are in form or body
		for key, expectedValue := range requestResponse.requestParams {
			jsonParamVal, _ := jsonParams[key]

			if r.FormValue(key) != expectedValue && jsonParamVal != expectedValue {
				h.notFound(w, requestResponse.returnHeaders)

				return
			}
		}

		h.setResponse(w, requestResponse.returnStatus, requestResponse.returnBody, requestResponse.returnHeaders)
	}))

	return h.server
}

// Close is using to close a fake server
func (h *HTTPFake) Close() {
	h.started = false

	h.server.Close()
}

func (h *HTTPFake) notFound(w http.ResponseWriter, headers map[string]string) {
	h.setResponse(w, 404, `{"error":"Not Found"}`, headers)
}

func (h *HTTPFake) setResponse(w http.ResponseWriter, status int, body string, headers map[string]string) {
	h.addHeaders(w, headers)
	w.WriteHeader(status)
	w.Write([]byte(body))
}

func (h *HTTPFake) addHeaders(w http.ResponseWriter, headers map[string]string) {
	for key, value := range headers {
		w.Header().Add(key, value)
	}
}

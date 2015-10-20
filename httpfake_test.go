package go_httpfake

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vladiacob/go_requester"
)

func TestAddRoute(t *testing.T) {
	httpFake := New()

	currentResponse1 := httpFake.AddRoute("test/query:approve", "GET", map[string]string{}, 200, `{"success": "true"}`, map[string]string{})
	currentResponse2 := httpFake.AddRoute("test/query:approve", "POST", map[string]string{}, 200, `{"success": "true"}`, map[string]string{})
	currentResponse3 := httpFake.AddRoute("test/query", "POST", map[string]string{"param1": "value1", "param2": "value2"}, 200, `{"success": "true"}`, map[string]string{"header1": "value1"})

	expected := map[string]map[string]RequestResponse{
		"test/query:approve": map[string]RequestResponse{
			"GET": RequestResponse{
				returnStatus:  200,
				returnBody:    `{"success": "true"}`,
				returnHeaders: map[string]string{},
				requestParams: map[string]string{},
			},
			"POST": RequestResponse{
				returnStatus:  200,
				returnBody:    `{"success": "true"}`,
				returnHeaders: map[string]string{},
				requestParams: map[string]string{},
			},
		},
		"test/query": map[string]RequestResponse{
			"POST": RequestResponse{
				returnStatus:  200,
				returnBody:    `{"success": "true"}`,
				returnHeaders: map[string]string{"header1": "value1"},

				requestParams: map[string]string{"param1": "value1", "param2": "value2"},
			},
		},
	}

	assert.Equal(t, expected, httpFake.routes)
	assert.Equal(t, true, currentResponse1)
	assert.Equal(t, true, currentResponse2)
	assert.Equal(t, true, currentResponse3)
}

func TestDelRoute(t *testing.T) {
	httpFake := New()

	httpFake.AddRoute("test/query:approve", "GET", map[string]string{}, 200, `{"success": "true"}`, map[string]string{})
	httpFake.AddRoute("test/query:approve", "POST", map[string]string{}, 200, `{"success": "true"}`, map[string]string{})
	httpFake.AddRoute("test/query", "POST", map[string]string{"param1": "value1", "param2": "value2"}, 200, `{"success": "true"}`, map[string]string{"header1": "value1"})

	currentResponse1 := httpFake.DelRoute("test/query:approve", "POST")
	currentResponse2 := httpFake.DelRoute("test/query", "POST")
	currentResponse3 := httpFake.DelRoute("test/wrong", "POST")

	expected := map[string]map[string]RequestResponse{
		"test/query:approve": map[string]RequestResponse{
			"GET": RequestResponse{
				returnStatus:  200,
				returnBody:    `{"success": "true"}`,
				returnHeaders: map[string]string{},
				requestParams: map[string]string{},
			},
		},
	}

	assert.Equal(t, expected, httpFake.routes)
	assert.Equal(t, true, currentResponse1)
	assert.Equal(t, true, currentResponse2)
	assert.Equal(t, false, currentResponse3)
}

func TestStart(t *testing.T) {
	requester := go_requester.New(http.DefaultClient)

	testCases := []struct {
		testName string

		route      string
		method     string
		parameters map[string]string
		status     int
		body       string
		headers    map[string]string

		expectedRoute      string
		expectedMethod     string
		expectedParameters map[string]string
		expectedStatus     int
		expectedResponse   string
	}{
		{
			testName: "Success Cases: GET without query parameters",

			route:      "test/query:approve",
			method:     "GET",
			parameters: map[string]string{},
			status:     200,
			body:       `{"success": "true"}`,
			headers:    map[string]string{},

			expectedRoute:      "test/query:approve",
			expectedMethod:     "GET",
			expectedParameters: map[string]string{},
			expectedStatus:     200,
			expectedResponse:   `{"success": "true"}`,
		},
		{
			testName: "Success Cases: GET with query parameters",

			route:      "test/query?query1=1&query2=2",
			method:     "GET",
			parameters: map[string]string{},
			status:     200,
			body:       `{"success": "true"}`,
			headers:    map[string]string{},

			expectedRoute:      "test/query?query1=1&query2=2",
			expectedMethod:     "GET",
			expectedParameters: map[string]string{"query1": "1", "query2": "2"},
			expectedStatus:     200,
			expectedResponse:   `{"success": "true"}`,
		},
		{
			testName: "Success Cases: POST without query parameters and body parameters",

			route:      "test/query:approve",
			method:     "POST",
			parameters: map[string]string{},
			status:     201,
			body:       `{"success": "true"}`,
			headers:    map[string]string{},

			expectedRoute:      "test/query:approve",
			expectedMethod:     "POST",
			expectedParameters: map[string]string{},
			expectedStatus:     201,
			expectedResponse:   `{"success": "true"}`,
		},
		{
			testName: "Success Cases: POST without query parameters and with body parameters",

			route:      "test/query:approve",
			method:     "POST",
			parameters: map[string]string{"param1": "value1", "param2": "value2"},
			status:     400,
			body:       `{"success": "true"}`,
			headers:    map[string]string{},

			expectedRoute:      "test/query:approve",
			expectedMethod:     "POST",
			expectedParameters: map[string]string{"param1": "value1", "param2": "value2"},
			expectedStatus:     400,
			expectedResponse:   `{"success": "true"}`,
		},
		{
			testName: "Success Cases: POST with query parameters and with body parameters",

			route:      "test/query?query1=1",
			method:     "POST",
			parameters: map[string]string{"param1": "value1", "param2": "value2"},
			status:     400,
			body:       `{"success": "true"}`,
			headers:    map[string]string{"query1": "1"},

			expectedRoute:      "test/query?query1=1",
			expectedMethod:     "POST",
			expectedParameters: map[string]string{"param1": "value1", "param2": "value2", "query1": "1"},
			expectedStatus:     400,
			expectedResponse:   `{"success": "true"}`,
		},
		{
			testName: "Success Cases: POST with query parameters and without body parameters",

			route:      "test/query?query1=1&query2=2",
			method:     "POST",
			parameters: map[string]string{},
			status:     500,
			body:       `{"success": "true"}`,
			headers:    map[string]string{},

			expectedRoute:      "test/query?query1=1&query2=2",
			expectedMethod:     "POST",
			expectedParameters: map[string]string{"query1": "1", "query2": "2"},
			expectedStatus:     500,
			expectedResponse:   `{"success": "true"}`,
		},
	}

	for _, test := range testCases {
		httpFake := New()
		httpFake.AddRoute(test.expectedRoute, test.expectedMethod, test.expectedParameters, test.expectedStatus, test.body, test.headers)

		fakeServer := httpFake.Start()

		var currentResponse string
		requester.Make(test.method, fmt.Sprintf("%s/%s", fakeServer.URL, test.route), test.parameters, &currentResponse)

		assert.Equal(t, test.expectedResponse, currentResponse)

		httpFake.Close()
	}
}

func TestStartError(t *testing.T) {
	requester := go_requester.New(http.DefaultClient)

	testCases := []struct {
		testName string

		route      string
		method     string
		parameters map[string]string
		status     int
		body       string
		headers    map[string]string

		expectedRoute      string
		expectedMethod     string
		expectedParameters map[string]string
		expectedStatus     int
		expectedResponse   string
	}{
		{
			testName: "Error Cases: Route not found",

			route:      "test/wrong",
			method:     "GET",
			parameters: map[string]string{},
			status:     404,
			body:       ``,
			headers:    map[string]string{},

			expectedRoute:      "test/query",
			expectedMethod:     "GET",
			expectedParameters: map[string]string{},
			expectedStatus:     200,
			expectedResponse:   `{"error":"Not Found"}`,
		},
		{
			testName: "Error Cases: Method not found",

			route:      "test/query",
			method:     "GET",
			parameters: map[string]string{},
			status:     404,
			body:       ``,
			headers:    map[string]string{},

			expectedRoute:      "test/query",
			expectedMethod:     "POST",
			expectedParameters: map[string]string{},
			expectedStatus:     200,
			expectedResponse:   `{"error":"Not Found"}`,
		},
		{
			testName: "Error Cases: Almost one param not found",

			route:      "test/query",
			method:     "GET",
			parameters: map[string]string{},
			status:     404,
			body:       ``,
			headers:    map[string]string{},

			expectedRoute:      "test/query",
			expectedMethod:     "GET",
			expectedParameters: map[string]string{"param1": "value1"},
			expectedStatus:     200,
			expectedResponse:   `{"error":"Not Found"}`,
		},
	}

	for _, test := range testCases {
		httpFake := New()
		httpFake.AddRoute(test.expectedRoute, test.expectedMethod, test.expectedParameters, test.expectedStatus, test.body, test.headers)

		fakeServer := httpFake.Start()

		var currentResponse string
		requester.Make(test.method, fmt.Sprintf("%s/%s", fakeServer.URL, test.route), test.parameters, &currentResponse)

		assert.Equal(t, test.expectedResponse, currentResponse)

		httpFake.Close()
	}
}

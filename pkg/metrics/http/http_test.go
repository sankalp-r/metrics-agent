package http

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpSource(t *testing.T) {
	testMetrics := "# test metric1\nmeteric_1_total 1\n# test metric2\nmetric_2_total 2"
	testStatus := `{"test-metric-1":1,"test-metric-2":2}`
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		if request.URL.String() == "/metrics" {
			w.Header().Set("Content-Type", textContentHeader)
			fmt.Fprintf(w, testMetrics)

		} else if request.URL.String() == "/v2/status" {
			w.Header().Set("Content-Type", jsonContentHeader)
			fmt.Fprintf(w, testStatus)

		}

	}))

	defer testServer.Close()
	testCases := []struct {
		desc             string
		endpoint         string
		expectedResponse string
		expectedErr      error
	}{
		{
			desc:             "fetch metrics",
			endpoint:         testServer.URL + "/metrics",
			expectedResponse: "# test metric1\nmeteric_1_total 1\n# test metric2\nmetric_2_total 2",
			expectedErr:      nil,
		},
		{
			desc:             "fetch status",
			endpoint:         testServer.URL + "/v2/status",
			expectedResponse: "# HELP test_metric_1 test metric 1\n# TYPE test_metric_1 gauge\ntest_metric_1 1\n# HELP test_metric_2 test metric 2\n# TYPE test_metric_2 gauge\ntest_metric_2 2\n",
			expectedErr:      nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			testHttpSource := NewHttpSource(test.endpoint, map[string]string{})
			response, err := testHttpSource.Collect()
			assert.Equal(t, test.expectedErr, err)
			assert.Equal(t, test.expectedResponse, string(response))

		})
	}

}

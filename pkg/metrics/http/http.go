package http

import (
	"bytes"
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"
	"github.com/sankalp-r/metrics-agent/pkg/metrics"
	"go.uber.org/zap"
	"io"
	"net/http"
	"reflect"
	"strings"
	"sync"
)

const (
	textContentHeader = "text/plain"
	jsonContentHeader = "application/json"
)

// httpSource
type httpSource struct {
	endpoint       string
	header         map[string]string
	counterMetrics map[string]prometheus.Counter
	gaugeMetrics   map[string]prometheus.Gauge
	lock           sync.RWMutex
	metricRegistry *prometheus.Registry
}

// NewHttpSource return new http-source
func NewHttpSource(endpoint string, headers map[string]string) metrics.Source {
	return &httpSource{
		endpoint:       endpoint,
		header:         headers,
		counterMetrics: make(map[string]prometheus.Counter),
		gaugeMetrics:   make(map[string]prometheus.Gauge),
		lock:           sync.RWMutex{},
		metricRegistry: prometheus.NewRegistry(),
	}
}

// Collect http-metrics
func (h *httpSource) Collect() ([]byte, error) {
	request, err := http.NewRequest("GET", h.endpoint, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range h.header {
		request.Header.Set(k, v)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil || response.StatusCode != http.StatusOK {
		return nil, err
	}

	defer response.Body.Close()
	contentType := response.Header.Get("Content-Type")
	switch contentType {
	case textContentHeader:
		return h.parseStringMetrics(response.Body), nil
	case jsonContentHeader:
		return h.parseJsonMetrics(response.Body), nil
	}
	return nil, nil
}

// parseStringMetrics
func (h *httpSource) parseStringMetrics(response io.Reader) []byte {
	b, err := io.ReadAll(response)
	if err != nil {
		zap.L().Error(err.Error())
		return nil
	}
	return b
}

// parseJsonMetrics
func (h *httpSource) parseJsonMetrics(response io.Reader) []byte {
	var buffer bytes.Buffer

	b, err := io.ReadAll(response)
	if err != nil {
		zap.L().Error(err.Error())
		return nil
	}

	var keys map[string]interface{}

	err = json.Unmarshal(b, &keys)
	if err != nil {
		zap.L().Error(err.Error())
		return nil
	}

	for k, v := range keys {
		if reflect.TypeOf(v).Kind() == reflect.Float64 {
			metricName := strings.ReplaceAll(k, "-", "_")
			metricHelp := strings.ReplaceAll(k, "-", " ")
			val := v.(float64)
			h.lock.Lock()
			if _, exist := h.gaugeMetrics[metricName]; !exist {
				h.gaugeMetrics[metricName] = prometheus.NewGauge(prometheus.GaugeOpts{
					Name: metricName,
					Help: metricHelp,
				})
				if err = h.metricRegistry.Register(h.gaugeMetrics[metricName]); err != nil {
					zap.L().Error(err.Error())
				}
			}
			h.gaugeMetrics[metricName].Set(val)
			h.lock.Unlock()
		}
	}

	encoder := expfmt.NewEncoder(&buffer, expfmt.FmtText)
	mf, err := h.metricRegistry.Gather()
	if err != nil {
		zap.L().Error(err.Error())
		return nil
	}
	for _, metric := range mf {
		err = encoder.Encode(metric)
		if err != nil {
			zap.L().Error(err.Error())
		}

	}
	return buffer.Bytes()
}

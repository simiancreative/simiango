package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/simiancreative/simiango/config"
	"github.com/simiancreative/simiango/logger"
)

func NewMockServiceConfig() *MockServiceConfig {
	return &MockServiceConfig{}
}

type MockResp map[string]any

type ConfigItem struct {
	Status int
	Value  string
	Resp   MockResp
}

type Metrics struct {
	requests int
}

type MockServiceConfig struct {
	urlKey      string
	urlSuffix   string
	configItems []ConfigItem
	server      *httptest.Server

	metrics Metrics
}

func (m *MockServiceConfig) SetConfigs(configSlice ...ConfigItem) *MockServiceConfig {
	m.configItems = configSlice

	return m
}

func (m *MockServiceConfig) SetURLKey(key string) *MockServiceConfig {
	m.urlKey = key

	return m
}

func (m *MockServiceConfig) SetURLSuffix(suffix string) *MockServiceConfig {
	m.urlSuffix = suffix

	return m
}

func (m *MockServiceConfig) StartServer() *httptest.Server {
	config.SetupTest()

	m.server = MockService(
		m.urlKey,
		m.urlSuffix,
		m.testHandler(),
	)

	return m.server
}

func (m *MockServiceConfig) HandledRequestsCount() int {
	return m.metrics.requests
}

func (m *MockServiceConfig) findConfigItem(value string) *ConfigItem {
	for _, config := range m.configItems {
		if config.Value == value {
			return &config
		}
	}

	return nil
}

func (m *MockServiceConfig) testHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m.metrics.requests++
		config := m.findConfigItem(r.URL.Path)
		logger.Warnf("Evaluating request: %v (%v)", r.URL.Path, config)

		w.Header().Set("Content-Type", "application/json")

		if config != nil {
			w.WriteHeader(config.Status)
			valBytes, _ := json.Marshal(config.Resp)
			w.Write(valBytes)
			return
		}

		logger.Warnf("Returning default response: (%v)", r.URL.Path)
		w.WriteHeader(204)
	}
}

func MockService(
	urlKey, urlSuffix string,
	handler func(w http.ResponseWriter, r *http.Request),
) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(handler))

	url := server.URL

	// add url suffix to the server url only if it is not empty
	if urlSuffix != "" {
		url = fmt.Sprintf("%v/%v", server.URL, urlSuffix)
	}

	os.Setenv(
		urlKey,
		fmt.Sprintf("%v", url),
	)

	return server
}

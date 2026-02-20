package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthCheck(t *testing.T) {

	app := &application{}

	req, err := http.NewRequest(http.MethodGet, "/v1/healthcheck", nil)

	require.NoError(t, err)

	rw := httptest.NewRecorder()

	app.healthCheck(rw, req)

	assert.Equal(t, http.StatusOK, rw.Code)

	var response struct {
		Message struct {
			Status  string `json:"status"`
			Version string `json:"version"`
		} `json:"message"`
	}

	err = json.Unmarshal(rw.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "available", response.Message.Status)
	assert.Equal(t, "1.0.0", response.Message.Version)

}

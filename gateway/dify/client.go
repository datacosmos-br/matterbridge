package dify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client talks to a Dify instance over its REST API.
type Client struct {
	BaseURL string
	APIKey  string
	HTTP    *http.Client
}

func New(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL: baseURL,
		APIKey:  apiKey,
		HTTP:    &http.Client{Timeout: 30 * time.Second},
	}
}

// WorkflowRequest is the payload sent to Dify workflow run endpoint.
type WorkflowRequest struct {
	Inputs  map[string]any `json:"inputs"`
	User    string         `json:"user"`
	Mode    string         `json:"response_mode"`
	TraceID string         `json:"trace_id,omitempty"`
}

// WorkflowResponse is the minimal Dify workflow response.
type WorkflowResponse struct {
	WorkflowRunID string         `json:"workflow_run_id"`
	Status        string         `json:"status"`
	Outputs       map[string]any `json:"data,omitempty"`
}

// RunWorkflow calls Dify workflow run endpoint with idempotency key.
func (c *Client) RunWorkflow(req WorkflowRequest, idempotencyKey string) (*WorkflowResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	httpReq, err := http.NewRequest("POST", c.BaseURL+"/v1/workflows/run", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")
	if idempotencyKey != "" {
		httpReq.Header.Set("X-Idempotency-Key", idempotencyKey)
	}
	resp, err := c.HTTP.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("dify workflow run failed: %s", resp.Status)
	}
	var out WorkflowResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

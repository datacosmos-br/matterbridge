package dify

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// maxErrorBody bounds how much of a failed response body is kept in errors.
const maxErrorBody = 512

// workflowRunPath is the Dify workflow execution endpoint, below DifyURL.
const workflowRunPath = "/v1/workflows/run"

// Workflow failures.
var (
	ErrWorkflowStatus   = errors.New("workflow request failed")
	ErrWorkflowResponse = errors.New("invalid workflow response")
)

// SummaryRequest is one summary request sent to the workflow.
type SummaryRequest struct {
	Channel        string // gateway the summary is requested for
	Transcript     string // recent messages, see Format
	RequestedBy    string // corporate user of the requester
	IdempotencyKey string // identity of the requesting message
}

// Client runs the summary workflow of a Dify application.
type Client struct {
	endpoint string
	apiKey   string
	http     *http.Client
}

// NewClient returns a client for the workflow configured in s.
func NewClient(s Settings) *Client {
	return &Client{
		endpoint: s.URL.JoinPath(workflowRunPath).String(),
		apiKey:   s.APIKey,
		http:     &http.Client{Timeout: s.WorkflowTimeout},
	}
}

type workflowRequest struct {
	Inputs       workflowInputs `json:"inputs"`
	ResponseMode string         `json:"response_mode"`
	User         string         `json:"user"`
}

type workflowInputs struct {
	Channel     string `json:"channel"`
	Transcript  string `json:"transcript"`
	RequestedBy string `json:"requested_by"`
}

type workflowResponse struct {
	Data *struct {
		Status  string `json:"status"`
		Error   string `json:"error"`
		Outputs *struct {
			Text *string `json:"text"`
		} `json:"outputs"`
	} `json:"data"`
}

// Summarize runs the workflow in blocking mode and returns its outputs.text.
func (c *Client) Summarize(ctx context.Context, req SummaryRequest) (string, error) {
	body, err := json.Marshal(workflowRequest{
		Inputs: workflowInputs{
			Channel:     req.Channel,
			Transcript:  req.Transcript,
			RequestedBy: req.RequestedBy,
		},
		ResponseMode: "blocking",
		User:         req.RequestedBy,
	})
	if err != nil {
		return "", fmt.Errorf("encode workflow request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("build workflow request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Idempotency-Key", req.IdempotencyKey)

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("run workflow: %w", err)
	}

	text, err := readOutput(resp)

	return text, errors.Join(err, resp.Body.Close())
}

// readOutput returns the outputs.text of a workflow run response.
func readOutput(resp *http.Response) (string, error) {
	if resp.StatusCode != http.StatusOK {
		snippet, err := io.ReadAll(io.LimitReader(resp.Body, maxErrorBody))
		if err != nil {
			return "", fmt.Errorf("%w: status %s (reading body: %w)", ErrWorkflowStatus, resp.Status, err)
		}

		return "", fmt.Errorf("%w: status %s: %s", ErrWorkflowStatus, resp.Status, snippet)
	}

	return decodeOutput(resp.Body)
}

func decodeOutput(r io.Reader) (string, error) {
	var out workflowResponse

	err := json.NewDecoder(r).Decode(&out)
	if err != nil {
		return "", fmt.Errorf("decode workflow response: %w", err)
	}

	switch {
	case out.Data == nil:
		return "", fmt.Errorf("%w: no data", ErrWorkflowResponse)
	case out.Data.Status != "succeeded":
		return "", fmt.Errorf("%w: run %q: %s", ErrWorkflowResponse, out.Data.Status, out.Data.Error)
	case out.Data.Outputs == nil || out.Data.Outputs.Text == nil:
		return "", fmt.Errorf("%w: no outputs.text", ErrWorkflowResponse)
	case *out.Data.Outputs.Text == "":
		return "", fmt.Errorf("%w: outputs.text is empty", ErrWorkflowResponse)
	}

	return *out.Data.Outputs.Text, nil
}

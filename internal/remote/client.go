package remote

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type RemoteClient struct {
	baseURL    string
	token      string
	httpClient *http.Client
	logger     *Logger
}

type RemoteOperation struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Status    string                 `json:"status"`
	Payload   map[string]interface{} `json:"payload"`
	StartTime time.Time              `json:"start_time"`
	EndTime   time.Time              `json:"end_time"`
}

func NewRemoteClient(baseURL string, token string, logger *Logger) *RemoteClient {
	return &RemoteClient{
		baseURL: baseURL,
		token:   token,
		httpClient: &http.Client{
			Timeout: time.Second * 30,
		},
		logger: logger,
	}
}

func (rc *RemoteClient) ExecuteOperation(opType string, payload map[string]interface{}) (*RemoteOperation, error) {
	operation := &RemoteOperation{
		Type:      opType,
		Payload:   payload,
		StartTime: time.Now(),
	}

	data, err := json.Marshal(operation)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal operation: %w", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/operations", rc.baseURL), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", rc.token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := rc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute remote operation: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("remote operation failed with status: %d", resp.StatusCode)
	}

	var result RemoteOperation
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

func (rc *RemoteClient) GetOperationStatus(operationID string) (*RemoteOperation, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/operations/%s", rc.baseURL, operationID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", rc.token))

	resp, err := rc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get operation status: %w", err)
	}
	defer resp.Body.Close()

	var operation RemoteOperation
	if err := json.NewDecoder(resp.Body).Decode(&operation); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &operation, nil
}

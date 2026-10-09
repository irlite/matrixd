package hub

import (
	"encoding/json"
	"errors"
)

func parseMessage(message string) (request, error) {
	var req request
	if err := json.Unmarshal([]byte(message), &req); err != nil {
		return request{}, err
	}
	if req.Type == "" || req.Action == "" {
		return req, errors.New("type and action required")
	}
	return req, nil
}

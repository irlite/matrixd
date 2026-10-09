package hub

import (
	"encoding/json"

	"github.com/irlite/matrixd/internal/wire"
)

type request struct {
	ID       string          `json:"id"`
	Type     string          `json:"type"`
	Action   string          `json:"action"`
	Content  json.RawMessage `json:"content"`
	ConnName string          `json:"-"`
}

type handler func(req request) wire.Outgoing

type errorContent struct {
	Message string `json:"message"`
}

func success(req request, content any) wire.Outgoing {
	return wire.Outgoing{
		ID:      req.ID,
		Type:    req.Type,
		Action:  req.Action,
		Status:  wire.StatusOK,
		Content: content,
	}
}

func fail(req request, message string) wire.Outgoing {
	return wire.Outgoing{
		ID:      req.ID,
		Type:    req.Type,
		Action:  req.Action,
		Status:  wire.StatusError,
		Content: errorContent{Message: message},
	}
}

package hub

import (
	"encoding/json"
	"log"

	"github.com/irlite/matrixd/internal/wire"
)

func marshal(out wire.Outgoing) string {
	data, err := json.Marshal(out)
	if err != nil {
		log.Printf("marshal %s.%s: %v", out.Type, out.Action, err)
		data, _ = json.Marshal(wire.Outgoing{
			ID:     out.ID,
			Type:   out.Type,
			Action: out.Action,
			Status: wire.StatusError,
			Content: errorContent{
				Message: "response could not be encoded",
			},
		})
	}
	return string(data)
}

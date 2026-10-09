package send

type TextContent struct {
	RoomID string `json:"room_id"`
	Body   string `json:"body"`
}

func Text(_ string, _ TextContent) (any, error) {
	return nil, nil
}

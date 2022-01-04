package dto

type MessageDto struct {
	ID    int    `json:"id"`
	Token string `json:"token"`
}

func NewMessageDto() MessageDto {
	return MessageDto{}
}

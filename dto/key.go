package dto

type KeyDto struct {
	Key string `json:"key"`
}

func NewKeyDto() KeyDto {
	return KeyDto{}
}

package usecase

type Event struct {
	Type string `json:"type" validate:"required"`
	Text string `json:"text" validate:"required"`
}

type Message struct {
	ID     string `json:"id" validate:"required"`
	Value  string `json:"value" validate:"required"`
	SendAt string `json:"sendAt" validate:"required"`
}

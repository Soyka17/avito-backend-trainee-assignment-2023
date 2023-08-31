package usecase

type ReportRequestDTO struct {
	Uid    int    `json:"user_id"`
	Slug   string `json:"slug"`
	After  string `json:"after"`
	Before string `json:"before"`
}

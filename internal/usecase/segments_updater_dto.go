package usecase

type SegmentsUpdaterDTO struct {
	UserId         int      `json:"user_id"`
	NewSegments    []string `json:"new"`
	RemoveSegments []string `json:"remove"`
}

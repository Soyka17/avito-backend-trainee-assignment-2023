package usecase

import (
	"AvitoInternship/internal/controllers/repository"
	"AvitoInternship/internal/entity"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
)

type SegmentDeleter struct {
	db     repository.Interface
	logger *zap.Logger
}

func NewSegmentDeleter(db repository.Interface, l *zap.Logger) *SegmentDeleter {
	return &SegmentDeleter{db: db, logger: l}
}

func (s *SegmentDeleter) Run(rawBody []byte) (int, map[string]any) {

	var delSeg entity.Segment
	err := json.Unmarshal(rawBody, &delSeg)
	if err != nil {
		e := "Unable to unmarshall json body"
		s.logger.Warn(e, zap.Error(err))
		return 400, map[string]any{"error": e}
	}

	slug := delSeg.Slug
	err = s.db.DeleteSegment(slug)
	if err != nil {
		switch err.(type) {
		case repository.SegmentNotFound:
			s.logger.Debug(err.Error(), zap.Error(err))
			return 400, map[string]any{"error": err.Error()}

		default:
			e := fmt.Sprintf("Received error with delete %s segment", slug)
			s.logger.Error(e, zap.Error(err))
			return 400, map[string]any{"error": e}
		}
	}

	return 200, map[string]any{"message": fmt.Sprintf("%s segment succesfully deleted", slug)}
}

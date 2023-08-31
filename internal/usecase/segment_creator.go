package usecase

import (
	"AvitoInternship/internal/controllers/repository"
	"AvitoInternship/internal/entity"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
)

type SegmentCreator struct {
	db     repository.Interface
	logger *zap.Logger
}

func NewSegmentCreator(db repository.Interface, l *zap.Logger) *SegmentCreator {
	return &SegmentCreator{db: db, logger: l}
}

func (s *SegmentCreator) Run(rawBody []byte) (int, map[string]any) {
	var newSeg entity.Segment
	err := json.Unmarshal(rawBody, &newSeg)
	if err != nil {
		e := "Unable to unmarshall json body"
		s.logger.Warn(e, zap.Error(err))
		return 400, map[string]any{"error": e}
	}

	err = s.db.CreateSegment(&newSeg)
	if err != nil {
		switch err.(type) {

		case repository.SegmentAlreadyExist:
			s.logger.Debug(err.Error(), zap.Error(err))
			return 400, map[string]any{"error": err.Error()}

		default:
			e := "Received error with creating new segment"
			s.logger.Error(e, zap.Error(err))
			return 400, map[string]any{"error": e}
		}
	}

	return 201, map[string]any{"message": fmt.Sprintf("New segment %s succesfully created", newSeg.Slug)}
}

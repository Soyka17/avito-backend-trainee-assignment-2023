package usecase

import (
	"AvitoInternship/internal/controllers/repository"
	"AvitoInternship/internal/entity"
	"encoding/json"
	"go.uber.org/zap"
)

type SegmentDeadlineSetter struct {
	db     repository.Interface
	logger *zap.Logger
}

func NewSegmentDeadlineSetter(db repository.Interface, l *zap.Logger) *SegmentDeadlineSetter {
	return &SegmentDeadlineSetter{db: db, logger: l}
}

func (s *SegmentDeadlineSetter) Run(rawBody []byte) (int, map[string]any) {
	var dd DeadlineDTO
	err := json.Unmarshal(rawBody, &dd)
	if err != nil {
		e := "Unable to unmarshall json body"
		s.logger.Warn(e, zap.Error(err))
		return 400, map[string]any{"error": e}
	}
	var seg *entity.Segment
	seg, err = s.db.GetUserActiveSegment(dd.Uid, dd.Slug)
	if err != nil {

		e := "Received error while get segment for set deadline"
		s.logger.Error(e, zap.Error(err), zap.Any("uid", dd.Uid), zap.Any("slug", dd.Slug))
		return 500, map[string]any{"error": e}
	}
	seg.EndDate = dd.Deadline.Time
	err = s.db.UnBindSegment(dd.Uid, *seg)
	if err != nil {
		e := "Received error while set deadline"
		s.logger.Error(e, zap.Error(err), zap.Any("uid", dd.Uid), zap.Any("slug", dd.Slug))
		return 500, map[string]any{"error": e}
	}

	return 201, map[string]any{"message": "Deadline successfully set"}
}

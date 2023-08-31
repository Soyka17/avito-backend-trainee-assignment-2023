package usecase

import (
	"AvitoInternship/internal/controllers/repository"
	"AvitoInternship/internal/entity"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type UserSegmentsChecker struct {
	db     repository.Interface
	logger *zap.Logger
}

func NewUserSegmentsChecker(db repository.Interface, l *zap.Logger) *UserSegmentsChecker {
	return &UserSegmentsChecker{db: db, logger: l}
}

func (u *UserSegmentsChecker) Run(req string) (int, map[string]any) {
	id, err := strconv.Atoi(req)
	if err != nil {
		return 404, map[string]any{"error": "Uncorrected user id"}
	}
	var segments []entity.Segment
	segments, err = u.db.GetUserActiveSegments(id)

	if err != nil {
		switch err.(type) {
		case repository.UserNotFound:
			u.logger.Debug(err.Error(), zap.Error(err))
			return 404, map[string]any{"error": err.Error()}
		default:
			e := "Received error with find user"
			u.logger.Error(e, zap.Error(err))
			return 500, map[string]any{"error": e}
		}
	}
	resp := make([]map[string]any, len(segments))
	for i := range segments {
		resp[i] = make(map[string]any)
		resp[i]["id"] = segments[i].Id
		resp[i]["slug"] = segments[i].Slug
		resp[i]["begin_date"] = segments[i].BeginDate.Format(time.RFC3339)
		if segments[i].EndDate == (time.Time{}) {
			resp[i]["end_date"] = ""
		} else {
			resp[i]["end_date"] = segments[i].EndDate.Format(time.RFC3339)
		}
	}
	return 200, map[string]any{"user_id": id, "active_segments": resp}
}

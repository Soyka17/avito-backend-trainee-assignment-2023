package usecase

import (
	"AvitoInternship/internal/controllers/report_repository"
	"AvitoInternship/internal/controllers/repository"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"time"
)

type ReportCreator struct {
	url        string
	db         repository.Interface
	logger     *zap.Logger
	repStorage report_repository.Interface
}

func NewReportCreator(url string, db repository.Interface, rs report_repository.Interface, l *zap.Logger) *ReportCreator {
	return &ReportCreator{url: url, db: db, repStorage: rs, logger: l}
}

func (s *ReportCreator) Run(rawBody []byte) (int, map[string]any) {
	var req ReportRequestDTO
	err := json.Unmarshal(rawBody, &req)
	if err != nil {
		e := "Unable to unmarshall json body"
		s.logger.Warn(e, zap.Error(err))
		return 400, map[string]any{"error": e}
	}
	var after, before time.Time
	after, err = time.Parse("2006-01-02 15:04:05", req.After)
	if err != nil {
		return 404, map[string]any{"message": fmt.Sprintf("After time uncorrect")}
	}
	before, err = time.Parse("2006-01-02 15:04:05", req.Before)
	if err != nil {
		return 404, map[string]any{"message": fmt.Sprintf("Before time uncorrect")}
	}
	segs, err := s.db.GetUserSegmentHistory(req.Uid, req.Slug, after, before)
	if err != nil {
		e := "Received error while get user segment history"
		s.logger.Error(e, zap.Error(err), zap.Any("req", req))
		return 500, map[string]any{"message": e}
	}

	path, err := s.repStorage.CreateReport(req.Uid, segs)
	if err != nil {
		e := "Received error while save user segment history"
		s.logger.Error(e, zap.Error(err), zap.Any("req", req))
		return 500, map[string]any{"message": e}
	}
	link := fmt.Sprintf("%s/%s", s.url, path)
	return 201, map[string]any{"report": link}
}

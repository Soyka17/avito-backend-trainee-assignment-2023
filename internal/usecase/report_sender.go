package usecase

import (
	"AvitoInternship/internal/controllers/repository"
	"fmt"
	"go.uber.org/zap"
	"os"
)

type ReportSender struct {
	url    string
	path   string
	db     repository.Interface
	logger *zap.Logger
}

func NewReportSender(url string, path string, db repository.Interface, l *zap.Logger) *ReportSender {
	return &ReportSender{url: url, path: path, db: db, logger: l}
}

func (s *ReportSender) Run(req string) (int, string) {
	path := fmt.Sprintf("%s/%s.csv", s.path, req)
	s.logger.Info(path)
	_, err := os.Open(path)
	if err != nil {
		e := "File not found"
		s.logger.Error(e, zap.Error(err), zap.Any("req", req))
		return 404, e
	}
	return 200, path
}

package report_repository

import (
	"AvitoInternship/internal/entity"
	"encoding/csv"
	"fmt"
	"go.uber.org/zap"
	"os"
	"strconv"
	"time"
)

type CsvReportMaker struct {
	logger *zap.Logger
	path   string
}

func NewReportMaker(l *zap.Logger, p string) *CsvReportMaker {
	return &CsvReportMaker{l, p}
}

func (m *CsvReportMaker) CreateReport(uid int, segs []entity.Segment) (string, error) {
	path := fmt.Sprintf("%s/%s_%s_%s", m.path, strconv.Itoa(uid), segs[0].Slug, time.Now().Format("2006-01-02_15:04:05"))
	err := m.saveCSV(uid, segs, path)
	if err != nil {
		return "", err
	}
	return path, nil
}

func (m *CsvReportMaker) saveCSV(uid int, segs []entity.Segment, path string) error {
	records := make([][]string, len(segs))
	for i := range segs {
		records[i] = append(records[i], strconv.Itoa(segs[i].Id))
		records[i] = append(records[i], segs[i].Slug)
		records[i] = append(records[i], segs[i].BeginDate.Format("2006-01-02 15:04:05"))
		end := ""
		if segs[i].EndDate != (time.Time{}) {
			end = segs[i].EndDate.Format("2006-01-02 15:04:05")
		}
		records[i] = append(records[i], end)
	}

	f, err := os.Create(fmt.Sprintf("%s.csv", path))
	defer f.Close()

	if err != nil {
		e := "Receive error while create report file"
		m.logger.Error(e, zap.Error(err), zap.Any("user", uid))
		return err
	}

	w := csv.NewWriter(f)
	err = w.WriteAll(records)

	if err != nil {
		e := "Receive error while save report file"
		m.logger.Error(e, zap.Error(err), zap.Any("user", uid))
		return err
	}

	return nil
}

package report_repository

import "AvitoInternship/internal/entity"

type Interface interface {
	CreateReport(uid int, segs []entity.Segment) (string, error)
}

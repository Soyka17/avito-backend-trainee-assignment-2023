package repository

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type PgSegmentDTO struct {
	Id        int              `json:"id"`
	Slug      string           `json:"slug"`
	BeginDate pgtype.Timestamp `json:"begin_date"`
	EndDate   pgtype.Timestamp `json:"end_date"`
}

package entity

import "time"

type Segment struct {
	Id        int       `json:"id"`
	Slug      string    `json:"slug"`
	BeginDate time.Time `json:"begin_date"`
	EndDate   time.Time `json:"end_date"`
}

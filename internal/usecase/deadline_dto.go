package usecase

import (
	"strings"
	"time"
)

type DeadlineDTO struct {
	Uid      int        `json:"user_id"`
	Slug     string     `json:"slug"`
	Deadline CustomDate `json:"deadline"`
}

type CustomDate struct {
	time.Time
}

func (c *CustomDate) UnmarshalJSON(b []byte) (err error) {
	layout := "2006-01-02 15:04:05"

	s := strings.Trim(string(b), "\"") // remove quotes
	if s == "null" {
		return
	}
	c.Time, err = time.Parse(layout, s)
	return
}

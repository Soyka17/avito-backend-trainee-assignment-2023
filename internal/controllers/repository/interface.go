package repository

import (
	"AvitoInternship/internal/entity"
	"time"
)

type Interface interface {
	CreateUser() (int, error)
	DeleteUser(id int) error

	CreateSegment(segment *entity.Segment) error
	GetSegmentId(slug string) (int, error)
	DeleteSegment(slug string) error

	BindSegment(uid int, segment entity.Segment) error
	UnBindSegment(uid int, segment entity.Segment) error
	GetUserSegments(id int) ([]entity.Segment, error)
	GetUserActiveSegment(uid int, slug string) (*entity.Segment, error)
	GetUserActiveSegments(id int) ([]entity.Segment, error)
	GetUserInactiveSegments(id int) ([]entity.Segment, error)
	GetUserSegmentHistory(uid int, slug string, after time.Time, before time.Time) ([]entity.Segment, error)
}

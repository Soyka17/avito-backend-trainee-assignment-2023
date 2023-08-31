package usecase

import (
	"AvitoInternship/internal/controllers/repository"
	"AvitoInternship/internal/entity"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"time"
)

type UserSegmentsUpdater struct {
	db     repository.Interface
	logger *zap.Logger
}

func NewUserSegmentsUpdater(db repository.Interface, l *zap.Logger) *UserSegmentsUpdater {
	return &UserSegmentsUpdater{db: db, logger: l}
}

//func (s *UserSegmentsUpdater) Run(rawBody []byte) (int, map[string]any) {
//
//	var update SegmentsUpdaterDTO
//	err := json.Unmarshal(rawBody, &update)
//	if err != nil {
//		s.logger.Warn("Unable to unmarshall json body", zap.Error(err))
//		return 400, map[string]any{"error": "Unable to unmarshall json body"}
//	}
//	var rawSegs []entity.Segment
//	rawSegs, err = s.db.GetUserSegments(update.UserId)
//	if err != nil {
//		switch err.(type) {
//		case repository.UserNotFound:
//			s.logger.Debug(err.Error(), zap.Error(err))
//			return 404, map[string]any{"error": err.Error()}
//		default:
//			e := "Received error with find user"
//			s.logger.Error(e, zap.Error(err))
//			return 500, map[string]any{"error": e}
//		}
//	}
//
//	segs := make(map[string]entity.Segment)
//	for _, i := range rawSegs {
//		segs[i.Slug] = i
//	}
//
//	var okBind []string
//	notBind := make(map[string]string)
//	for _, seg := range update.NewSegments {
//		if _, ok := segs[seg]; !ok {
//			err = s.db.BindSegment(update.UserId, seg)
//			if err != nil {
//				switch err.(type) {
//				case repository.SegmentNotExist:
//					s.logger.Debug(err.Error(), zap.Error(err))
//					notBind[seg] = err.Error()
//				default:
//					e := "Received error with bind segment"
//					s.logger.Error(e, zap.Error(err))
//					notBind[seg] = err.Error()
//				}
//			} else {
//				okBind = append(okBind, seg)
//			}
//		} else {
//			notBind[seg] = "Segment already binded"
//		}
//	}
//	var okUnBind []string
//	notUnBind := make(map[string]string)
//	for _, seg := range update.RemoveSegments {
//		err = s.db.UnBindSegment(update.UserId, seg)
//		if err != nil {
//			s.logger.Debug(err.Error(), zap.Error(err))
//			notUnBind[seg] = err.Error()
//		} else {
//			okUnBind = append(okUnBind, seg)
//		}
//	}
//
//	resp := make(map[string]any)
//	resp["bind"] = okBind
//	resp["not_bind"] = notBind
//	resp["unbind"] = okUnBind
//	resp["not_unbind"] = notUnBind
//	return 200, resp
//}

func (s *UserSegmentsUpdater) Run(rawBody []byte) (int, map[string]any) {
	var update SegmentsUpdaterDTO
	err := json.Unmarshal(rawBody, &update)
	if err != nil {
		s.logger.Warn("Unable to unmarshall json body", zap.Error(err))
		return 400, map[string]any{"error": "Unable to unmarshall json body"}
	}
	onSegsRaw, err := s.db.GetUserActiveSegments(update.UserId)

	onSegs := make(map[string]entity.Segment, len(onSegsRaw))
	for i, seg := range onSegsRaw {
		onSegs[seg.Slug] = onSegsRaw[i]
	}

	okUnBind, notUnBind := s.unbindSegments(onSegs, update)
	okBind, notBind := s.bindSegments(onSegs, update)

	resp := make(map[string]any)
	resp["ok_bind"] = okBind
	resp["not_bind"] = notBind
	resp["ok_unbind"] = okUnBind
	resp["not_unbind"] = notUnBind
	if len(notBind) != 0 || len(notUnBind) != 0 {
		return 400, resp
	}
	return 200, resp
}

func (s *UserSegmentsUpdater) unbindSegments(onSegs map[string]entity.Segment, update SegmentsUpdaterDTO) ([]string, map[string]string) {
	var okUnBind []string
	notUnBind := make(map[string]string)
	for _, removeS := range update.RemoveSegments {
		if _, ok := onSegs[removeS]; ok {
			err := s.db.UnBindSegment(update.UserId, onSegs[removeS])
			if err != nil {
				e := "Received error while unbind segment"
				s.logger.Warn(e, zap.Error(err), zap.Any("user", update.UserId), zap.Any("segment", removeS))
				notUnBind[removeS] = e
				continue
			}
			okUnBind = append(okUnBind, removeS)
		} else {
			e := "The user is missing a segment or the segment is already unbinded"
			s.logger.Debug(e, zap.Any("user", update.UserId), zap.Any("segment", removeS))
			notUnBind[removeS] = e
		}
	}
	return okUnBind, notUnBind
}
func (s *UserSegmentsUpdater) bindSegments(onSegs map[string]entity.Segment, update SegmentsUpdaterDTO) ([]string, map[string]string) {
	var okBind []string
	notBind := make(map[string]string)
	for _, bindS := range update.NewSegments {
		if _, ok := onSegs[bindS]; ok {
			notBind[bindS] = "Segment already binded"
			continue
		}
		id, err := s.db.GetSegmentId(bindS)
		if err != nil {
			e := fmt.Sprintf("Unable to find id by slug=%s", bindS)
			s.logger.Debug(e, zap.Error(err))
			notBind[bindS] = err.Error()
			continue
		}
		newSeg := entity.Segment{Id: id, Slug: bindS, BeginDate: time.Now()}
		err = s.db.BindSegment(update.UserId, newSeg)

		if err != nil {
			switch err.(type) {
			case repository.SegmentAlreadyExist:
				s.logger.Debug(err.Error(), zap.Error(err))
			default:
				e := "Received error while bind segment"
				s.logger.Error(e, zap.Error(err), zap.Any("user", update.UserId), zap.Any("segment", bindS))
			}
			notBind[bindS] = err.Error()
			continue
		}

		okBind = append(okBind, bindS)
	}
	return okBind, notBind
}

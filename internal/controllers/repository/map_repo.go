package repository

//
//import (
//	"AvitoInternship/internal/entity"
//	"time"
//)
//
//type MapRepository struct {
//	users    map[int]*entity.User
//	segments map[string]*entity.Segment
//}
//
//func NewMapRepository() *MapRepository {
//	return &MapRepository{users: make(map[int]*entity.User), segments: make(map[string]*entity.Segment)}
//}
//
//func (r *MapRepository) GetUser(id int) (*entity.User, error) {
//	if user, ok := r.users[id]; ok {
//		return user, nil
//	}
//	return nil, UserNotFound{id: id}
//}
//func (r *MapRepository) CreateUser() (int, error) {
//	newUser := &entity.User{Id: r.getLastUserId(), Segments: make(map[string]*entity.UserSegment)}
//	r.users[newUser.Id] = newUser
//	return newUser.Id, nil
//}
//func (r *MapRepository) DeleteUser(id int) error {
//	if _, ok := r.users[id]; ok {
//		r.users[id] = nil
//		return nil
//	}
//	return UserNotFound{id: id}
//}
//
//func (r *MapRepository) CreateSegment(segment *entity.Segment) error {
//	segment.Id = r.getLastSegmentId()
//	if _, ok := r.segments[segment.Slug]; ok {
//		return SegmentAlreadyExist{segment.Slug}
//	}
//	r.segments[segment.Slug] = segment
//	return nil
//}
//func (r *MapRepository) DeleteSegment(slug string) error {
//	if _, ok := r.segments[slug]; !ok {
//		return SegmentNotFound{slug}
//	}
//	for i := range r.users {
//		for _, s := range r.users[i].Segments {
//			if s.Reference.Slug == slug {
//				s.EndDate = time.Now()
//			}
//		}
//	}
//	if _, ok := r.segments[slug]; ok {
//		delete(r.segments, slug)
//	}
//	return nil
//}
//
//func (r *MapRepository) getLastUserId() int {
//	max := 0
//	for i := range r.users {
//		if r.users[i].Id > max {
//			max = r.users[i].Id
//		}
//	}
//	return max + 1
//}
//func (r *MapRepository) getLastSegmentId() int {
//	max := 0
//	for i := range r.segments {
//		if r.segments[i].Id > max {
//			max = r.segments[i].Id
//		}
//	}
//	return max + 1
//}
//
//func (r *MapRepository) BindSegment(id int, slug string) error {
//	user, ok := r.users[id]
//	if !ok {
//		return UserNotFound{}
//	}
//
//	reference, refExist := r.segments[slug]
//	if !refExist {
//		return SegmentNotExist{slug: slug}
//	}
//	newSeg := &entity.UserSegment{Reference: reference, BeginDate: time.Now()}
//	r.users[user.Id].Segments[slug] = newSeg
//	return nil
//}
//func (r *MapRepository) UnBindSegment(id int, slug string) error {
//	user, ok := r.users[id]
//	if !ok {
//		return UserNotFound{id}
//	}
//
//	if _, ok = r.segments[slug]; !ok {
//		return SegmentNotFound{slug}
//	}
//
//	if _, ok = r.users[id].Segments[slug]; !ok {
//		return UserMissingSegment{id, slug}
//	}
//
//	end := r.users[id].Segments[slug].EndDate
//	if (end != time.Time{}) && time.Now().After(end) {
//		return UserSegmentAlreadyUnbind{id, slug}
//	}
//
//	user.Segments[slug].EndDate = time.Now()
//
//	return nil
//}
//func (r *MapRepository) GetUserSegments(id int) ([]*entity.UserSegment, error) {
//	user, err := r.GetUser(id)
//	if err != nil {
//		return nil, err
//	}
//	var res []*entity.UserSegment
//	for _, v := range user.Segments {
//		res = append(res, v)
//	}
//	return res, nil
//}
//func (r *MapRepository) GetUserActiveSegmentsJSON(id int) (map[string]any, error) {
//	segs, err := r.GetUserSegments(id)
//	if err != nil {
//		return nil, err
//	}
//	var res []map[string]any
//	for i := range segs {
//		if segs[i].EndDate == (time.Time{}) {
//			curr := make(map[string]any)
//			curr["slug"] = segs[i].Reference.Slug
//			curr["begin"] = segs[i].BeginDate.Format(time.RFC3339)
//			curr["end"] = ""
//			res = append(res, curr)
//		}
//	}
//	return map[string]any{"segments": res}, nil
//}

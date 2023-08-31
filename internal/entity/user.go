package entity

type User struct {
	Id       int
	Segments map[string]*Segment
}

package repository

import (
	"fmt"
)

type UserNotFound struct {
	id int
}

func (e UserNotFound) Error() string {
	return fmt.Sprintf("User with id = %d not found", e.id)
}

type UserAlreadyExist struct {
	id int
}

func (e UserAlreadyExist) Error() string {
	return fmt.Sprintf("User with id = %d already exist", e.id)
}

type SegmentNotFound struct {
	slug string
}

func (e SegmentNotFound) Error() string {
	return fmt.Sprintf("%s segment not found", e.slug)
}

type SegmentAlreadyExist struct {
	slug string
}

func (e SegmentAlreadyExist) Error() string {
	return fmt.Sprintf("%s segment already exist", e.slug)
}

type SegmentNotExist struct {
	slug string
}

func (e SegmentNotExist) Error() string {
	return fmt.Sprintf("%s segment not exist", e.slug)
}

type UserMissingSegment struct {
	id   int
	slug string
}

func (e UserMissingSegment) Error() string {
	return fmt.Sprintf("User with id = %d didn't have %s segment", e.id, e.slug)
}

type UserSegmentAlreadyUnbind struct {
	id   int
	slug string
}

func (e UserSegmentAlreadyUnbind) Error() string {
	return fmt.Sprintf("%s segment for user with id = %d already unbinded", e.slug, e.id)
}

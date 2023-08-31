package report_repository

import "fmt"

type UserNotFound struct {
	id int
}

func (e UserNotFound) Error() string {
	return fmt.Sprintf("User with id = %d not found", e.id)
}

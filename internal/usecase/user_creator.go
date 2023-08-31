package usecase

import (
	"AvitoInternship/internal/controllers/repository"
	"fmt"
	"go.uber.org/zap"
)

type UserCreator struct {
	db     repository.Interface
	logger *zap.Logger
}

func NewUserCreator(db repository.Interface, l *zap.Logger) *UserCreator {
	return &UserCreator{db: db, logger: l}
}

func (u *UserCreator) Run(req string) (int, map[string]any) {
	id, err := u.db.CreateUser()
	resp := make(map[string]any)
	if err != nil {
		e := "Received error with creating new user"
		u.logger.Error(e, zap.Error(err))
		resp["message"] = e
		return 500, resp
	}
	resp["message"] = fmt.Sprintf("New user with id=%d succesfully created", id)
	return 201, resp
}

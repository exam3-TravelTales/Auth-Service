package service

import (
	pb "auth/genproto/users"
	"auth/pkg/logger"
	"auth/storage/postgres"
	"database/sql"
	"log/slog"
)

type UserService struct {
	pb.UnimplementedUserServer
	Repo   *postgres.UserRepo
	Logger *slog.Logger
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{
		Repo:   postgres.NewUserRepository(db),
		Logger: logger.NewLogger(),
	}
}

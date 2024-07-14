package handler

import (
	"auth/genproto/users"
	"log/slog"
)

type Handler struct {
	User users.UserClient
	Log  *slog.Logger
}

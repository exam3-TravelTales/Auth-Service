package service

import (
	pb "auth/genproto/users"
	"auth/pkg/logger"
	"auth/storage/postgres"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
)

type UserService struct {
	pb.UnimplementedUserServer
	Repo *postgres.UserRepo
	Log  *slog.Logger
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{
		Repo: postgres.NewUserRepository(db),
		Log:  logger.NewLogger(),
	}
}

func (u *UserService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	u.Log.Info("Register rpc method started")
	res, err := u.Repo.CreateUser(ctx, req)
	if err != nil {
		u.Log.Error(err.Error())
		return nil, err
	}
	u.Log.Info("Register rpc method finished")
	return res, nil
}
func (u *UserService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.UserInfo, error) {
	u.Log.Info("Login rpc method started")
	res, err := u.Repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		u.Log.Error(err.Error())
		return nil, err
	}
	if res.Password != req.Password {
		u.Log.Error("Password is incorrect")
		return nil, fmt.Errorf("password is incorrect")
	}

	u.Log.Info("Login rpc method finished")
	return res, nil
}

func (u *UserService) GetProfile(ctx context.Context, id *pb.UserId) (*pb.GetProfileResponse, error) {
	u.Log.Info("GetProfile rpc method started")
	res, err := u.Repo.GetUserProfile(ctx, id)
	if err != nil {
		u.Log.Error(err.Error())
		return nil, err
	}
	u.Log.Info("GetProfile rpc method finished")
	return res, nil
}

func (u *UserService) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	u.Log.Info("UpdateProfile rpc method started")
	res, err := u.Repo.UpdateUser(ctx, req)
	if err != nil {
		u.Log.Error(err.Error())
		return nil, err
	}
	u.Log.Info("UpdateProfile rpc method finished")
	return res, nil
}

func (u *UserService) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	u.Log.Info("GetUsers rpc method started")
	res, err := u.Repo.GetUsers(ctx, req)
	if err != nil {
		u.Log.Error(err.Error())
		return nil, err
	}
	u.Log.Info("GetUsers rpc method finished")
	return res, nil
}

func (u *UserService) DeleteUser(ctx context.Context, req *pb.UserId) (*pb.BoolResponse, error) {
	u.Log.Info("DeleteUser rpc method started")
	err := u.Repo.DeleteUser(ctx, req.Id)
	if err != nil {
		u.Log.Error(err.Error())
		return &pb.BoolResponse{Success: false}, err
	}
	u.Log.Info("DeleteUser rpc method finished")
	return &pb.BoolResponse{Success: true}, nil
}

func (u *UserService) EmailRecovery(ctx context.Context, req *pb.EmailRecoveryRequest) (*pb.BoolResponse, error) {
	u.Log.Info("EmailRecovery rpc method started")
	user, err := u.Repo.GetUserByID(ctx, req.UserId)
	if err != nil {
		u.Log.Error(err.Error())
		return &pb.BoolResponse{Success: false}, err
	}
	if user.Password != req.OldPassword {
		u.Log.Error("Password is incorrect")
		return &pb.BoolResponse{Success: false}, errors.New("password is incorrect")
	}

	err = u.Repo.UpdatePassword(ctx, req)
	if err != nil {
		u.Log.Error(err.Error())
		return &pb.BoolResponse{Success: false}, err
	}
	u.Log.Info("EmailRecovery rpc method finished")
	return &pb.BoolResponse{Success: true}, nil
}

func (u *UserService) Activity(ctx context.Context, req *pb.UserId) (*pb.ActivityResponse, error) {
	u.Log.Info("Activity rpc method started")
	res, err := u.Repo.GetUserActivity(ctx, req.Id)
	if err != nil {
		u.Log.Error(err.Error())
		return nil, err
	}
	u.Log.Info("Activity rpc method finished")
	return res, nil
}

func (u *UserService) Follow(ctx context.Context, req *pb.FollowRequest) (*pb.FollowResponse, error) {
	u.Log.Info("Follow rpc method started")

	res, err := u.Repo.Follow(ctx, req.FollowerId, req.FollowingId)
	if err != nil {
		u.Log.Error(err.Error())
		return nil, err
	}
	u.Log.Info("Follow rpc method finished")
	return res, nil
}

func (u *UserService) Followers(ctx context.Context, req *pb.FollowersRequest) (*pb.FollowersResponse, error) {
	u.Log.Info("Followers rpc method started")
	res, err := u.Repo.GetFollowers(ctx, req.UserId, req.Limit, req.Offset)
	if err != nil {
		u.Log.Error(err.Error())
	}
	u.Log.Info("Followers rpc method finished")
	return res, nil
}

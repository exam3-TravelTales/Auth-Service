package postgres

import (
	pb "auth/genproto/users"
	"context"
	"fmt"
	"reflect"
	"testing"
)

func TestCreateUser(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	con := NewUserRepository(db)
	res, err := con.CreateUser(context.Background(), &pb.RegisterRequest{
		Email:    "test@test.com",
		Password: "test",
		Username: "test",
		FullName: "test",
	})
	if err != nil {
		fmt.Println(err)
	}
	req := &pb.RegisterResponse{
		Id:        "dfb52830-c101-4114-bd07-97a94cce70ad",
		Email:     "test@test.com",
		Username:  "test",
		FullName:  "test",
		CreatedAt: "2024-07-16T16:19:28.613112+05:00",
	}
	if !reflect.DeepEqual(res, req) {
		t.Errorf("CreateUser returned %+v, want %+v", res, req)
	}
}

func TestGetUserByID(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	con := NewUserRepository(db)
	res, err := con.GetUserByID(context.Background(), "dfb52830-c101-4114-bd07-97a94cce70ad")
	if err != nil {
		fmt.Println(err)
	}
	req := &pb.UserInfo{
		Id:               "dfb52830-c101-4114-bd07-97a94cce70ad",
		Email:            "test@test.com",
		Password:         "test",
		Username:         "test",
		FullName:         "test",
		CountriesVisited: 0,
	}
	if !reflect.DeepEqual(res, req) {
		t.Errorf("GetUserByID returned %+v, want %+v", res, req)
	}
}

func TestGetUserProfile(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	con := NewUserRepository(db)
	res, err := con.GetUserProfile(context.Background(), &pb.UserId{Id: "dfb52830-c101-4114-bd07-97a94cce70ad"})
	if err != nil {
		fmt.Println(err)
	}
	req := &pb.GetProfileResponse{
		Id:               "dfb52830-c101-4114-bd07-97a94cce70ad",
		Email:            "test@test.com",
		Username:         "test",
		FullName:         "test",
		CountriesVisited: 0,
		CreatedAt:        "2024-07-16T16:19:28.613112+05:00",
		UpdatedAt:        "2024-07-16T16:19:28.613112+05:00",
	}
	if !reflect.DeepEqual(res, req) {
		t.Errorf("GetUserProfile returned %+v, want %+v", res, req)
	}
}

func TestGetUserByEmail(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	con := NewUserRepository(db)
	res, err := con.GetUserByEmail(context.Background(), "test@test.com")
	if err != nil {
		fmt.Println(err)
	}
	req := &pb.UserInfo{
		Id:               "dfb52830-c101-4114-bd07-97a94cce70ad",
		Email:            "test@test.com",
		Password:         "test",
		Username:         "test",
		FullName:         "test",
		CountriesVisited: 0,
	}
	if !reflect.DeepEqual(res, req) {
		t.Errorf("GetUserByID returned %+v, want %+v", res, req)
	}
}

func TestUpdateUser(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	con := NewUserRepository(db)
	res, err := con.UpdateUser(context.Background(), &pb.UpdateProfileRequest{
		Id:               "dfb52830-c101-4114-bd07-97a94cce70ad",
		FullName:         "test",
		CountriesVisited: 0,
		Bio:              "nma gap",
	})
	if err != nil {
		fmt.Println(err)
	}
	req := &pb.UpdateProfileResponse{
		Id:               "dfb52830-c101-4114-bd07-97a94cce70ad",
		Email:            "test@test.com",
		Username:         "test",
		FullName:         "test",
		CountriesVisited: 0,
		Bio:              "nma gap",
	}
	if !reflect.DeepEqual(res, req) {
		t.Errorf("UpdateUser returned %+v, want %+v", res, req)
	}
}

func TestDeleteUser(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	con := NewUserRepository(db)
	err = con.DeleteUser(context.Background(), "dfb52830-c101-4114-bd07-97a94cce70ad")
	if err != nil {
		fmt.Println(err)
	}
}

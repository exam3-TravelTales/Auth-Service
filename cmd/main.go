package main

import (
	"auth/api"
	"auth/api/handler"
	"auth/genproto/users"
	"auth/pkg/logger"
	"auth/service"
	"auth/storage/postgres"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
)

func main() {
	db, err := postgres.ConnectDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	fmt.Println("Starting server...")
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("error while listening: %v", err)
	}
	defer lis.Close()
	userService := service.NewUserService(db)
	server := grpc.NewServer()
	users.RegisterUserServer(server, userService)
	log.Printf("server listening at %v", lis.Addr())

	go func() {
		err = server.Serve(lis)
		if err != nil {
			log.Fatalf("error while serving: %v", err)
		}
	}()

	hand := NewHandler()
	router := api.Router(hand)
	log.Println("server is running")
	log.Fatal(router.Run(":8085"))

}
func NewHandler() *handler.Handler {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Panic(err)
	}
	return &handler.Handler{User: users.NewUserClient(conn), Log: logger.NewLogger()}
}

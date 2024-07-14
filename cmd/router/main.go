package main

import (
	"auth/api"
	"auth/api/handler"
	"auth/genproto/users"
	"auth/pkg/logger"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
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

package main

import (
	"log"
	"net"

	"user-service/internal/database"
	"user-service/internal/service"

	pb "github.com/Moldirkab/ap2_final_user_service_generated"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	database.ConnectDB()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterUserServiceServer(
		grpcServer,
		&service.AuthService{},
	)

	reflection.Register(grpcServer)

	log.Println("User Service running on :50051")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}

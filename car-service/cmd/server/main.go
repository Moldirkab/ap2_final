package main

import (
	"car_service/internal/cache"
	"car_service/internal/delivery/grpc"
	"car_service/internal/messaging"
	"car_service/internal/repository"
	"car_service/internal/usecase"
	"log"
	"net"

	carpb "github.com/Moldirkab/ap2_final_car_service_generated/carpb"

	grpcServer "google.golang.org/grpc"
)

func main() {
	// PostgreSQL connection
	db, err := repository.NewPostgresDB()
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	// Layers
	carRepo := repository.NewCarRepository(db)
	redisCache := cache.NewRedisCache()
	carUsecase := usecase.NewCarUsecase(carRepo, redisCache)
	messaging.StartNATSListener(carUsecase)
	carHandler := grpc.NewCarHandler(carUsecase)

	// gRPC server
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpcServer.NewServer()

	carpb.RegisterCarServiceServer(server, carHandler)

	log.Println("Car Service running on port 50052")

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

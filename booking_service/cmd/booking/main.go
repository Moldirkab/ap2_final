package main

import (
	"booking-service/internal/repository"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"booking-service/internal/clients"
	"booking-service/internal/events"
	grpcdelivery "booking-service/internal/transport/grpc"
	"booking-service/internal/usecase"

	bookingpb "github.com/Moldirkab/ap2_final_generated/bookingpb"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	_ = godotenv.Load()

	db, err := connectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userClient, err := clients.NewUserClient(os.Getenv("USER_SERVICE_ADDR"))
	if err != nil {
		log.Fatal(err)
	}

	carClient, err := clients.NewCarClient(os.Getenv("CAR_SERVICE_ADDR"))
	if err != nil {
		log.Fatal(err)
	}

	publisher, err := events.NewNatsPublisher(os.Getenv("NATS_URL"))
	if err != nil {
		log.Fatal(err)
	}

	bookingRepo := repository.NewBookingRepository(db)
	bookingUsecase := usecase.NewBookingUsecase(bookingRepo, userClient, carClient, publisher)
	bookingHandler := grpcdelivery.NewBookingHandler(bookingUsecase)

	port := os.Getenv("BOOKING_GRPC_PORT")

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()

	bookingpb.RegisterBookingServiceServer(server, bookingHandler)

	log.Println("Booking Service running on port", port)

	if err := server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}

func connectDB() (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

package main

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	carpb "github.com/Moldirkab/ap2_final_car_service_generated/carpb"
	bookingpb "github.com/Moldirkab/ap2_final_generated/bookingpb"
	userpb "github.com/Moldirkab/ap2_final_user_service_generated"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	carConn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer carConn.Close()

	userConn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer userConn.Close()

	bookingConn, err := grpc.NewClient("localhost:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer bookingConn.Close()

	carClient := carpb.NewCarServiceClient(carConn)
	userClient := userpb.NewUserServiceClient(userConn)
	bookingClient := bookingpb.NewBookingServiceClient(bookingConn)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))

	// ==================== USER ROUTES ====================

	router.POST("/auth/register", func(c *gin.Context) {
		var body struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		res, err := userClient.Register(context.Background(), &userpb.RegisterRequest{
			Email:    body.Email,
			Password: body.Password,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": res.Token})
	})

	router.POST("/auth/login", func(c *gin.Context) {
		var body struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		res, err := userClient.Login(context.Background(), &userpb.LoginRequest{
			Email:    body.Email,
			Password: body.Password,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": res.Token})
	})

	router.GET("/auth/profile", func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		res, err := userClient.ValidateToken(context.Background(), &userpb.TokenRequest{
			Token: token,
		})
		if err != nil || !res.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user_id": res.UserId,
			"role":    res.Role,
		})
	})

	// ==================== CAR ROUTES ====================

	router.POST("/cars", func(c *gin.Context) {
		type CreateCarHTTP struct {
			Brand       string  `json:"brand"`
			Model       string  `json:"model"`
			Year        int32   `json:"year"`
			PlateNumber string  `json:"plate_number"`
			PricePerDay float64 `json:"price_per_day"`
			Status      string  `json:"status"`
			Photo       string  `json:"photo"`
		}

		var body CreateCarHTTP
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		req := &carpb.CreateCarRequest{
			Brand:       body.Brand,
			Model:       body.Model,
			Year:        body.Year,
			PlateNumber: body.PlateNumber,
			PricePerDay: body.PricePerDay,
			Photo:       body.Photo,
		}

		res, err := carClient.CreateCar(context.Background(), req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, res.Car)
	})

	router.GET("/cars/:id", func(c *gin.Context) {
		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

		res, err := carClient.GetCar(context.Background(), &carpb.GetCarRequest{Id: id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, res.Car)
	})

	router.GET("/cars", func(c *gin.Context) {
		status := c.Query("status")
		brand := c.Query("brand")
		maxPriceStr := c.Query("max_price")

		var maxPrice float64
		if maxPriceStr != "" {
			maxPrice, _ = strconv.ParseFloat(maxPriceStr, 64)
		}

		res, err := carClient.ListCars(context.Background(), &carpb.ListCarsRequest{
			Brand:    brand,
			Status:   status,
			MaxPrice: maxPrice,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, res.Cars)
	})

	router.PUT("/cars/:id", func(c *gin.Context) {
		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

		var body struct {
			Brand       string  `json:"brand"`
			Model       string  `json:"model"`
			Year        int32   `json:"year"`
			PlateNumber string  `json:"plate_number"`
			PricePerDay float64 `json:"price_per_day"`
			Status      string  `json:"status"`
			Photo       string  `json:"photo"`
		}

		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		req := &carpb.UpdateCarRequest{
			Id:          id,
			Brand:       body.Brand,
			Model:       body.Model,
			Year:        body.Year,
			PlateNumber: body.PlateNumber,
			PricePerDay: body.PricePerDay,
			Status:      body.Status,
			Photo:       body.Photo,
		}

		res, err := carClient.UpdateCar(context.Background(), req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, res.Car)
	})

	router.DELETE("/cars/:id", func(c *gin.Context) {
		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

		res, err := carClient.DeleteCar(context.Background(), &carpb.DeleteCarRequest{Id: id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, res)
	})

	router.GET("/cars/:id/availability", func(c *gin.Context) {
		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

		res, err := carClient.CheckAvailability(context.Background(), &carpb.CheckAvailabilityRequest{CarId: id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"available": res.Available,
			"status":    res.Status,
		})
	})

	// ==================== BOOKING ROUTES ====================

	router.POST("/bookings", func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		var body struct {
			CarID     string `json:"car_id"`
			StartDate string `json:"start_date"`
			EndDate   string `json:"end_date"`
		}

		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		res, err := bookingClient.CreateBooking(context.Background(), &bookingpb.CreateBookingRequest{
			Token:     token,
			CarId:     body.CarID,
			StartDate: body.StartDate,
			EndDate:   body.EndDate,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, res)
	})

	router.GET("/bookings", func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		res, err := bookingClient.ListUserBookings(context.Background(), &bookingpb.ListUserBookingsRequest{
			Token: token,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, res.Bookings)
	})

	router.POST("/bookings/:id/confirm", func(c *gin.Context) {
		bookingID := c.Param("id")

		res, err := bookingClient.ConfirmBooking(
			context.Background(),
			&bookingpb.ConfirmBookingRequest{
				BookingId: bookingID,
			},
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, res)
	})

	router.POST("/bookings/:id/cancel", func(c *gin.Context) {
		bookingID := c.Param("id")

		res, err := bookingClient.CancelBooking(
			context.Background(),
			&bookingpb.CancelBookingRequest{
				BookingId: bookingID,
			},
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, res)
	})

	router.POST("/bookings/:id/complete", func(c *gin.Context) {
		bookingID := c.Param("id")

		res, err := bookingClient.CompleteBooking(
			context.Background(),
			&bookingpb.CompleteBookingRequest{
				BookingId: bookingID,
			},
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, res)
	})

	router.Run(":8080")
}

func extractToken(c *gin.Context) string {
	auth := c.GetHeader("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}
	return auth
}

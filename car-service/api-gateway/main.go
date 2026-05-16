package main

import (
	"context"
	"net/http"
	"strconv"

	carpb "github.com/Moldirkab/ap2_final_car_service_generated/carpb"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {

	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := carpb.NewCarServiceClient(conn)

	router := gin.Default()

	// CREATE CAR
	router.POST("/cars", func(c *gin.Context) {

		type CreateCarHTTP struct {
			Brand       string  `json:"brand"`
			Model       string  `json:"model"`
			Year        int32   `json:"year"`
			PlateNumber string  `json:"plate_number"`
			PricePerDay float64 `json:"price_per_day"`
			Status      string  `json:"status"`
		}

		var body CreateCarHTTP

		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		req := &carpb.CreateCarRequest{
			Brand:       body.Brand,
			Model:       body.Model,
			Year:        body.Year,
			PlateNumber: body.PlateNumber,
			PricePerDay: body.PricePerDay,
		}

		res, err := client.CreateCar(context.Background(), req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, res.Car)
	})

	// GET CAR
	router.GET("/cars/:id", func(c *gin.Context) {

		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

		res, err := client.GetCar(context.Background(), &carpb.GetCarRequest{
			Id: id,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, res.Car)
	})

	// LIST CARS
	router.GET("/cars", func(c *gin.Context) {

		status := c.Query("status")
		brand := c.Query("brand")

		maxPriceStr := c.Query("max_price")

		var maxPrice float64

		if maxPriceStr != "" {
			maxPrice, _ = strconv.ParseFloat(maxPriceStr, 64)
		}

		res, err := client.ListCars(context.Background(), &carpb.ListCarsRequest{
			Brand:    brand,
			Status:   status,
			MaxPrice: maxPrice,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, res.Cars)
	})

	// UPDATE CAR
	router.PUT("/cars/:id", func(c *gin.Context) {

		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

		var req carpb.UpdateCarRequest

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		req.Id = id

		res, err := client.UpdateCar(context.Background(), &req)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, res.Car)
	})

	// DELETE CAR
	router.DELETE("/cars/:id", func(c *gin.Context) {

		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

		res, err := client.DeleteCar(context.Background(), &carpb.DeleteCarRequest{
			Id: id,
		})

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

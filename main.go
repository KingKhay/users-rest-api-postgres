package main

import "github.com/gin-gonic/gin"

func main() {

	router := gin.Default()

	router.GET("/users")

	router.GET("/users/:id")

	router.POST("/users")

	router.PUT("/users/:id")

	router.DELETE("/users/:id")

	router.Run(":9300")
}

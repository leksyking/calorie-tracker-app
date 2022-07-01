package main

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/leksyking/calorie-tracker-app/routes"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	router := gin.New()
	//middlewares
	router.Use(gin.Logger(), cors.Default())

	//routes
	router.GET("/entries", routes.GetEntries)
	entryRoute := router.Group("/entry")
	{
		entryRoute.POST("/create", routes.AddEntry)
		entryRoute.GET("/:id", routes.GetEntryById)
		entryRoute.PUT("/update/:id", routes.UpdateEntry)
		entryRoute.DELETE("/delete/:id", routes.DeleteEntry)
	}
	ingredientsRoute := router.Group("/ingredients")
	{
		ingredientsRoute.GET("/:ingredients", routes.GetEntriesByIngredient)
		ingredientsRoute.PUT("/update/:id", routes.UpdateIngredient)
	}

	router.Run(":" + port)
}

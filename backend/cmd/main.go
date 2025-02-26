package main

import (
	"log"

	"github.com/bigelle/online-shop/backend/config"
	"github.com/bigelle/online-shop/backend/internal/database"
	"github.com/bigelle/online-shop/backend/internal/server"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("can't load .env file: %s", err.Error())
	}

	conf := config.New()

	db, err := database.Connect(conf.DatabaseDSN)
	if err != nil {
		log.Fatalf("error connecting to database: %s", err.Error())
	}
	defer database.Close(db)
	if err := database.Migrate(db); err != nil {
		log.Fatalf("can't automigrate: %s", err.Error())
	}

	r := gin.Default()

	//middleware? idk

	server.SetupRoutes(r, db)

	//FIXME
	r.Run(":8080")
}

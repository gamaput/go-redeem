package main

import (
	"github.com/gamaput/go-redeem/model"
	"github.com/gamaput/go-redeem/route"
)

func main() {
	db, _ := model.DBConnection()
	route.SetupRoutes(db)
}

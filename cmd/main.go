package main

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/my-Sakura/SMS/controller"
	"github.com/my-Sakura/SMS/utils"
)

func main() {
	r := gin.Default()

	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3307)/mysql")
	if err != nil {
		panic(err)
	}

	s := controller.NewSMSController(db)
	r.Use(utils.Cors())
	s.RegistRouter(r.Group("/api/v1"))

	r.Run(":8000")
}

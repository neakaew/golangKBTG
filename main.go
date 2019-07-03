package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	db1 "dome/school/database"
	"dome/school/schooldb"
)

func getPort() string {
	var port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
		fmt.Println("No Port In Heroku" + port)
	}
	return ":" + port
}

func main() {
	db, err := db1.Connect()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	r := gin.Default()
	r.GET("/api/todos", schooldb.GetTodos)
	r.GET("/api/todoGetByID/:id", schooldb.GetTodosByIdHandler)
	r.POST("/api/todoPost", schooldb.PostTodos)
	r.DELETE("/api/todoDeleteByID/:id", schooldb.DeleteTodosByIdHandler)
	r.Run(getPort())
}

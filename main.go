package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"dome/school/schooldb"
	_ "dome/school/schooldb"
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
	r := gin.Default()

	r.GET("/", schooldb.GetTodos)
	r.GET("/api/todoGetByID/:id", schooldb.GetTodosByIdHandler)
	r.POST("/api/todoPost", schooldb.PostTodos)
	r.DELETE("/api/todoDeleteByID/:id", schooldb.DeleteTodosByIdHandler)
	r.Run(":3344")
}

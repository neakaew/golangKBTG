package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"dome/school/schooldb"
	_ "dome/school/schooldb"
)

func main() {
	r := gin.Default()

	r.GET("/api/todoGet", schooldb.GetTodos)
	r.GET("/api/todoGetByID/:id", schooldb.GetTodosByIdHandler)
	r.POST("/api/todoPost", schooldb.PostTodos)
	r.DELETE("/api/todoDeleteByID/:id", schooldb.DeleteTodosByIdHandler)
	r.Run(":6655")
}

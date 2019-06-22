package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

func getTodos(c *gin.Context) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	stmt, err := db.Prepare("SELECT id, title, status FROM todos ORDER BY id ASC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rows, err := stmt.Query()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	todos := []Todo{}

	for rows.Next() {
		t := Todo{}

		err := rows.Scan(&t.ID, &t.Title, &t.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		todos = append(todos, t)
	}
	c.JSON(http.StatusOK, todos)
	defer db.Close()
}

func getTodosByIdHandler(c *gin.Context) {
	idParam := c.Param("id") // รับ paramiter id
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	db, _ := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	stmt, _ := db.Prepare("SELECT id, title, status FROM todos WHERE id=$1")

	row := stmt.QueryRow(id)
	t := Todo{}

	err2 := row.Scan(&t.ID, &t.Title, &t.Status)
	if err2 != nil {
		log.Fatal("error..", err2.Error())
	}

	c.JSON(http.StatusOK, t)
	defer db.Close()
}

func postTodos(c *gin.Context) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	t := &Todo{}
	if err := c.BindJSON(t); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var id int
	query := "INSERT INTO todos (title, status) VALUES ($1, $2) RETURNING id"
	row := db.QueryRow(query, &t.Title, &t.Status)
	err = row.Scan(&id)
	if err != nil {
		log.Fatal(err.Error(), id)
	}

	c.JSON(http.StatusOK, gin.H{
		"title":  t.Title,
		"status": t.Status,
	})
	defer db.Close()
}

func deleteTodosByIdHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	stmt, err := db.Prepare("DELETE FROM todos WHERE id=$1")
	rs, err := stmt.Exec(id)
	if err != nil {
		log.Fatalln(err)
	}
	c.JSON(http.StatusOK, rs)
	defer db.Close()
}

func main() {
	r := gin.Default()

	r.GET("/api/todos", getTodos)
	r.GET("/api/todos/:id", getTodosByIdHandler)
	r.POST("/api/todos", postTodos)
	r.DELETE("/api/seleceTodos/:id", deleteTodosByIdHandler)
	r.Run(":3434")
}

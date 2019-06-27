package schooldb

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"

	model "dome/school/models"
)

type todoes []model.Todo

func GetTodos(c *gin.Context) {
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

	todos := []model.Todo{}

	for rows.Next() {
		t := model.Todo{}

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

func GetTodosByIdHandler(c *gin.Context) {
	idParam := c.Param("id") // รับ paramiter id
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	db, _ := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	stmt, _ := db.Prepare("SELECT id, title, status FROM todos WHERE id=$1")

	row := stmt.QueryRow(id)
	t := model.Todo{}

	err2 := row.Scan(&t.ID, &t.Title, &t.Status)
	if err2 != nil {
		log.Fatal("error..", err2.Error())
	}

	c.JSON(http.StatusOK, t)
	defer db.Close()
}

func PostTodos(c *gin.Context) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	t := &model.Todo{}
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

func DeleteTodosByIdHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}
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

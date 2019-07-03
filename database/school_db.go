package school_db

import (
	"database/sql"
	"log"
	"os"

	model "dome/school/models"
)

type todoes []model.Todo

func Connect() (*sql.DB, error) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	return db, err
}

func main() {

}

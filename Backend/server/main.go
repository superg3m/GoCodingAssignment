package main

import (
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func GetDSN(dbEngine, host string, port int, user, password, dbname string) string {
	switch dbEngine {
	case "postgres":
		return fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname,
		)
	case "mysql":
		return fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?parseTime=true",
			user, password, host, port, dbname,
		)
	case "sqlite3":
		return dbname // For SQLite, `dbname` is the file path
	case "sqlserver":
		return fmt.Sprintf(
			"sqlserver://%s:%s@%s:%d?database=%s",
			user, password, host, port, dbname,
		)
	default:
		return ""
	}
}

// Note(Jovanni):
// If I was doing something more serious you would use some type of .env to hide these constants
const (
	DB_ENGINE = "mysql"
	HOST      = "localhost"
	PORT      = 3306
	USER      = "root"
	PASSWORD  = "P@55word"
	DBNAME    = "stoic"
)

type User struct {
	ID    int64  `db:"user_id" json:"user_id"`
	Name  string `db:"name" json:"name"`
	Email string `db:"email" json:"email"`
}

func main() {
	// Connect to DB
	dsn := GetDSN(DB_ENGINE, HOST, PORT, USER, PASSWORD, DBNAME)
	db, err := sqlx.Connect(DB_ENGINE, dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	e := echo.New()

	// GET /users/:id
	e.GET("/users/:id", func(c echo.Context) error {
		id := c.Param("id")

		var user User
		err := db.Get(&user, "SELECT id, name, email FROM users WHERE id=$1", id)
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
		}

		return c.JSON(http.StatusOK, user)
	})

	err = e.Start(":8080")
	if err != nil {
		return
	}
}

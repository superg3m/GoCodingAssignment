package main

import (
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/superg3m/server/Utility"
)

// Note(Jovanni):
// If I was doing something more serious you would use some type of
// .env to hide these constants
const (
	DB_ENGINE = "mysql"
	HOST      = "localhost"
	PORT      = 3306
	USER      = "root"
	PASSWORD  = "P@55word"
	DBNAME    = "stoic"
)

type User struct {
	ID         int64  `db:"user_id" json:"user_id"`
	Username   string `db:"user_name" json:"user_name"`
	Firstname  string `db:"first_name" json:"first_name"`
	Lastname   string `db:"last_name" json:"last_name"`
	Email      string `db:"email" json:"email"`
	UserStatus string `db:"user_status" json:"user_status"`
	Department string `db:"department" json:"department"`
}

func main() {
	// Connect to DB
	dsn := Utility.GetDSN(DB_ENGINE, HOST, PORT, USER, PASSWORD, DBNAME)
	db, err := sqlx.Connect(DB_ENGINE, dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
			return next(c)
		}
	})

	// Make sure to check to make sure user_name is unique
	e.POST("/User/Create", func(c echo.Context) error {
		var u User
		err := c.Bind(&u)
		if err != nil {
			return c.JSON(http.StatusNotFound, err.Error())
		}

		// do a CanCreate() check here

		sql := `INSERT INTO User (
        	user_name, 
        	first_name,
        	last_name,
            email,
            user_status,
            department
        ) VALUES (?, ?, ?, ?, ?, ?)
		`

		_, err = db.Exec(sql, u.Username, u.Firstname, u.Lastname, u.Email, u.UserStatus, u.Department)
		if err != nil {
			return c.JSON(http.StatusNotFound, err.Error())
		}

		return c.JSON(http.StatusOK, "User Created Successfully")
	})

	// GET /users/:id
	e.GET("/User/Get/:id", func(c echo.Context) error {
		id := c.Param("id")

		var user User
		err := db.Get(&user, "SELECT * FROM User WHERE user_id = ?", id)
		if err != nil {
			return c.JSON(http.StatusNotFound, err.Error())
		}

		return c.JSON(http.StatusOK, user)
	})

	err = e.Start(":8080")
	if err != nil {
		return
	}
}

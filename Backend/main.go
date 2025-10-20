package main

import (
	"github.com/labstack/echo/v4/middleware"
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

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{echo.GET, echo.POST, echo.PATCH, echo.DELETE, echo.OPTIONS},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
			return next(c)
		}
	})

	e.POST("/User/Create", func(c echo.Context) error {
		var u User
		err := c.Bind(&u)
		if err != nil {
			return c.JSON(http.StatusNotFound, err.Error())
		}

		_, err = db.Exec("SELECT * FROM User WHERE user_name = ?", u.Username)
		if err != nil {
			return c.JSON(http.StatusNotFound, "Found duplicate Username")
		}

		if !Utility.ValidEmail(u.Email) {
			return c.JSON(http.StatusNotFound, "Invalid Email")
		}

		sql := `INSERT INTO User (
        	user_name, 
        	first_name,
        	last_name,
            email,
            user_status,
            department
        ) VALUES (?, ?, ?, ?, ?, ?)
		`

		result, err := db.Exec(sql, u.Username, u.Firstname, u.Lastname, u.Email, u.UserStatus, u.Department)
		if err != nil {
			return c.JSON(http.StatusNotFound, err.Error())
		}

		u.ID, _ = result.LastInsertId()
		return c.JSON(http.StatusOK, u)
	})

	e.GET("/User/Get/All", func(c echo.Context) error {
		var users []User
		err := db.Select(&users, "SELECT * FROM User")
		if err != nil {
			return c.JSON(http.StatusNotFound, err.Error())
		}

		return c.JSON(http.StatusOK, users)
	})

	// NOTE(Jovanni):
	// This error can't really happen in a small scoped thing like this
	// However if there was like multiple people reading and writing to the
	// database then you could probably have a situation.
	e.GET("/User/Get/:id", func(c echo.Context) error {
		id := c.Param("id")

		var u User
		err := db.Get(&u, "SELECT * FROM User WHERE user_id = ?", id)
		if u.ID != 0 && err != nil {
			return c.JSON(http.StatusNotFound, err.Error())
		}

		return c.JSON(http.StatusOK, u)
	})

	// Make sure to check to make sure user_name is unique
	e.PATCH("/User/Update", func(c echo.Context) error {
		var u User
		err := c.Bind(&u)
		if u.ID != 0 && err != nil {
			return c.JSON(http.StatusNotFound, err.Error())
		}

		exists := false
		err = db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM User WHERE user_id = ?)", u.ID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		if !exists {
			return c.JSON(http.StatusConflict, "User ID doesn't exist exists")
		}

		if !Utility.ValidEmail(u.Email) {
			return c.JSON(http.StatusNotFound, "Invalid Email")
		}

		var dbUser User
		err = db.Get(&dbUser, "SELECT * FROM User WHERE user_name = ?", u.Username)
		if err == nil && dbUser.ID != u.ID {
			return c.JSON(http.StatusConflict, "Username already exists")
		}

		sql := `UPDATE User
		SET user_name = ?, 
		first_name = ?,
		last_name = ?,
		email = ?,
		user_status = ?,
		department = ?
		WHERE user_id = ?
		`

		_, err = db.Exec(sql, u.Username, u.Firstname, u.Lastname, u.Email, u.UserStatus, u.Department, u.ID)
		if err != nil {
			return c.JSON(http.StatusNotFound, err.Error())
		}

		return c.JSON(http.StatusOK, "User Updated Successfully")
	})

	// Make sure to check to make sure user_name is unique
	e.DELETE("/User/Delete", func(c echo.Context) error {
		var u User
		err := c.Bind(&u)
		if u.ID != 0 && err != nil {
			return c.JSON(http.StatusNotFound, err.Error())
		}

		sql := "DELETE FROM User WHERE user_id = ?"

		_, err = db.Exec(sql, u.ID)
		if err != nil {
			return c.JSON(http.StatusNotFound, err.Error())
		}

		return c.JSON(http.StatusOK, "User Deleted Successfully")
	})

	err = e.Start(":8080")
	if err != nil {
		return
	}
}

package main

import (
	"fmt"
	"github.com/labstack/echo/v4/middleware"
	"github.com/superg3m/server/Model"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/superg3m/server/Utility"
	_ "github.com/superg3m/server/docs"
)

// Note(Jovanni):
// Probably you could use dependency injection to not have this global
// but I think its fine for now.
var db *sqlx.DB

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

// @title Stoic API
// @version 1.0
// @description User management API for Stoic.
// @host localhost:8080
// @BasePath /
func main() {
	dsn := Utility.GetDSN(DB_ENGINE, HOST, PORT, USER, PASSWORD, DBNAME)
	var err error
	if db, err = sqlx.Connect(DB_ENGINE, dsn); err != nil {
		panic(err)
	}
	defer db.Close()

	e := echo.New()
	e.GET("/swagger/*", echoSwagger.WrapHandler)

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

	e.POST("/User/Create", CreateUser)
	e.GET("/User/Get/All", GetAllUsers)
	e.GET("/User/Get/:id", GetUserByID)
	e.PATCH("/User/Update", UpdateUser)
	e.DELETE("/User/Delete", DeleteUser)

	if err := e.Start(":8080"); err != nil {
		fmt.Println(err)
	}
}

// @Summary Create user
// @Description Creates a new user
// @Tags User
// @Accept json
// @Produce json
// @Param user body Model.User true "User data"
// @Success 200 {object} Model.User
// @Failure 400 {string} string
// @Failure 409 {string} string
// @Router /User/Create [post]
func CreateUser(c echo.Context) error {
	var u Model.User
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if !Utility.ValidEmail(u.Email) {
		return c.JSON(http.StatusBadRequest, "Invalid Email")
	}

	var dbUser Model.User
	if err := db.Get(&dbUser, "SELECT * FROM User WHERE user_name = ?", u.Username); err == nil {
		return c.JSON(http.StatusConflict, "Username already exists")
	}

	sqlStmt := `INSERT INTO User (user_name, first_name, last_name, email, user_status, department)
	            VALUES (?, ?, ?, ?, ?, ?)`
	result, err := db.Exec(sqlStmt, u.Username, u.Firstname, u.Lastname, u.Email, u.UserStatus, u.Department)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	u.ID, _ = result.LastInsertId()
	return c.JSON(http.StatusOK, u)
}

// @Summary Get all users
// @Tags User
// @Produce json
// @Success 200 {array} Model.User
// @Failure 404 {string} string
// @Router /User/Get/All [get]
func GetAllUsers(c echo.Context) error {
	var users []Model.User
	if err := db.Select(&users, "SELECT * FROM User"); err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, users)
}

// @Summary Get user by ID
// @Tags User
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} Model.User
// @Failure 404 {string} string
// @Router /User/Get/{id} [get]
func GetUserByID(c echo.Context) error {
	id := c.Param("id")
	var u Model.User
	if err := db.Get(&u, "SELECT * FROM User WHERE user_id = ?", id); err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, u)
}

// @Summary Update user
// @Tags User
// @Accept json
// @Produce json
// @Param user body Model.User true "Updated user data"
// @Success 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Router /User/Update [patch]
func UpdateUser(c echo.Context) error {
	var u Model.User
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if !Utility.ValidEmail(u.Email) {
		return c.JSON(http.StatusBadRequest, "Invalid Email")
	}

	var exists bool
	if err := db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM User WHERE user_id = ?)", u.ID); err != nil || !exists {
		return c.JSON(http.StatusBadRequest, "User ID doesn't exist")
	}

	var dbUser Model.User
	if err := db.Get(&dbUser, "SELECT * FROM User WHERE user_name = ?", u.Username); err == nil && dbUser.ID != u.ID {
		return c.JSON(http.StatusConflict, "Username already exists")
	}

	sqlStmt := `UPDATE User SET user_name=?, first_name=?, last_name=?, email=?, user_status=?, department=? WHERE user_id=?`
	if _, err := db.Exec(sqlStmt, u.Username, u.Firstname, u.Lastname, u.Email, u.UserStatus, u.Department, u.ID); err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, "User Updated Successfully")
}

// @Summary Delete user
// @Tags User
// @Accept json
// @Produce json
// @Param user body Model.User true "User data"
// @Success 200 {string} string
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /User/Delete [delete]
func DeleteUser(c echo.Context) error {
	var u Model.User
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	result, err := db.Exec("DELETE FROM User WHERE user_id = ?", u.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return c.JSON(http.StatusBadRequest, "User not found")
	}

	return c.JSON(http.StatusOK, "User Deleted Successfully")
}

package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"strconv"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func setupDB() (*sql.DB, error) {
	dbUser := os.Getenv("dbUser")
	dbPassword := os.Getenv("dbPassword")
	dbName := os.Getenv("dbName")

	connStr := fmt.Sprintf("postgresql://%s:%s@localhost:5436/%s?sslmode=disable", dbUser, dbPassword, dbName)

	//Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {

	db, err := setupDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := gin.Default()

	//Handling getting all Users
	router.GET("/users", func(c *gin.Context) {

		var listOfUsers []User
		rows, err := db.Query("SELECT * FROM users")
		defer rows.Close()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		for rows.Next() {

			var user User

			//Scan takes each row column and assigns the value to the user struct
			if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			listOfUsers = append(listOfUsers, user)
		}

		if len(listOfUsers) == 0 {
			c.IndentedJSON(http.StatusOK, []User{})
			return
		}

		c.IndentedJSON(http.StatusOK, listOfUsers)
	})

	//Handling creating new User
	router.POST("/users", func(c *gin.Context) {
		var user User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		//Execute the Insert query
		_, err := db.Exec("INSERT INTO users (name, email) VALUES ($1, $2)", user.Name, user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.IndentedJSON(http.StatusCreated, &user)
	})

	//Handling getting a single user by Id
	router.GET("/users/:id", func(c *gin.Context) {
		userIdStr := c.Param("id")
		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user id"})
			return
		}

		var user User
		if err := db.QueryRow("SELECT * FROM users WHERE id = $1", userId).Scan(&user.ID, &user.Name, &user.Email); err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "no user found with id"})
			return
		}

		c.JSON(http.StatusOK, &user)
	})

	//Handling updating User
	router.PUT("/users/:id", func(c *gin.Context) {
		var user User
		userIdStr := c.Param("id")
		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid user id",
				"status":  http.StatusBadRequest,
			})
			return
		}

		err = c.BindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		_, err = db.Exec("UPDATE users SET name = $1, email = $2 WHERE id = $3", user.Name, user.Email, userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.IndentedJSON(http.StatusOK, &user)
	})

	//Handling deleting users
	router.DELETE("/users/:id", func(c *gin.Context) {
		userIdStr := c.Param("id")
		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user id"})
			return
		}

		result, err := db.Exec("DELETE FROM users WHERE id = $1", userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "could not delete user"})
			return
		}

		affectedRows, _ := result.RowsAffected()
		if affectedRows == 0 {
			c.Status(http.StatusNoContent)
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
	})

	router.Run("localhost:9300")
}

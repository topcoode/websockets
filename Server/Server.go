package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
)

type User struct {
	Name     string //`json: "name"`
	Username string //`json: "username"`
	Status   string //`json: "status"`
	Message  string //`json: "message"`
	Online   string // `json: "online"`
}

// var clients = make(map[*websocket.Conn]bool)
var upgrader = websocket.Upgrader{}

func main() {
	router := gin.Default()
	router.POST("/client1", Client1)
	router.GET("/client2", Client2)
	if err := router.Run(":8080"); err != nil {
		fmt.Println("err on the port ...")
	}
}
func Client1(c *gin.Context) {
	var userdata User
	if err := c.ShouldBindJSON(&userdata); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(userdata)
	message := userdata.Message

	fmt.Println("the data of the message is :", message)

	onlineoffline := userdata.Status
	//offline := userdata.Offline
	fmt.Println(onlineoffline)

	if onlineoffline == "A" {
		fmt.Println("online")
	} else if onlineoffline == "B" {
		fmt.Println("offline")
	} else {
		fmt.Println("Given data is error")
	}
	c.JSON(http.StatusOK, gin.H{"message": "Data received successfully"})

	connStr := "postgres://postgres:123456@localhost/postgres?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
	}

	defer db.Close()
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully connected to the database!")
	dbdata := `INSERT INTO users(username,status) VALUES ($1,$2)`
	_, err = db.Exec(dbdata, userdata.Username, userdata.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}
	//sender--------------->

}

func Client2(c *gin.Context) {
	connStr := "postgres://postgres:123456@localhost/postgres?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
	}

	defer db.Close()
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully connected to the database!")
	rows, err := db.Query("SELECT username, status FROM users")
	if err != nil {
		fmt.Println(err)
		//c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Username, &user.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, user)
	}
	c.JSON(http.StatusOK, users)

}

package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

func dbFactory() (*sql.DB, error) {
	DB_NAME := os.Getenv("DB_NAME")
	DB_USER := os.Getenv("DB_USER")
	DB_PASS := os.Getenv("DB_PASS")

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASS, DB_NAME)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return db, errors.New("could not connect to database")
	}

	return db, nil
}

func getAlbums(c *gin.Context) {
	var albums []Album

	db, err := dbFactory()

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"statusCode": http.StatusInternalServerError, "error": err.Error()})
		return
	}

	defer db.Close()

	cur, err := db.Query("SELECT * FROM album")

	if err != nil {
		panic(err.Error())
	}

	defer cur.Close()

	for cur.Next() {
		var album Album

		if err := cur.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"statusCode": http.StatusInternalServerError, "error": "something went wrong"})
		}

		albums = append(albums, album)
	}

	c.IndentedJSON(http.StatusOK, albums)
}

func postAlbums(c *gin.Context) {
	var newAlbum Album

	if err := c.BindJSON(&newAlbum); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "error": "Could not add album. Try again"})
		return
	}

	db, err := dbFactory()

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"statusCode": http.StatusInternalServerError, "error": err.Error()})
		return
	}

	defer db.Close()

	_, err = db.Exec("INSERT INTO album (title, artist, price) VALUES ($1, $2, $3)", newAlbum.Title, newAlbum.Artist, newAlbum.Price)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"statusCode": http.StatusCreated, "message": "Album created successfully"})
}

func getAlbumByID(c *gin.Context) {
	id := c.Param("id")
	var album Album

	db, err := dbFactory()

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"statusCode": http.StatusInternalServerError, "error": err.Error()})
		return
	}

	defer db.Close()

	row := db.QueryRow("SELECT * FROM album WHERE id = $1", id)

	if err := row.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound, gin.H{"statusCode": http.StatusNotFound, "error": "Album not found"})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"statusCode": http.StatusInternalServerError, "error": err.Error()})
        return
	}

	c.IndentedJSON(http.StatusOK, album)
}

func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.POST("/albums", postAlbums)
	router.GET("/albums/:id", getAlbumByID)

	router.Run("localhost:8080")
}

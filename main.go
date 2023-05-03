package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

type Request struct {
	ID      int64           `json:"id"`
	Content json.RawMessage `json:"content"`
}

var db *sql.DB

func main() {
	// Connexion à la base de données PostgreSQL
	connStr := os.Getenv("POSTGRESQL_ADDON_URI")
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Création de la table pour stocker les requêtes
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS requests (id SERIAL PRIMARY KEY, content JSONB)`)
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/requests", handleRequest)
	e.GET("/requests", getRequests)
	e.GET("/requests/:id", getRequestByID)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8182"
	}

	log.Fatal(e.Start(":" + port))
}

func handleRequest(c echo.Context) error {
	if c.Request().Header.Get("Content-Type") != "application/json" {
		return c.String(http.StatusBadRequest, "Content-Type must be application/json")
	}

	var content json.RawMessage
	err := json.NewDecoder(c.Request().Body).Decode(&content)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid JSON")
	}

	result, err := db.Exec(`INSERT INTO requests (content) VALUES ($1)`, content)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to save request to database")
	}

	id, _ := result.LastInsertId()
	return c.JSON(http.StatusCreated, map[string]int64{"id": id})
}

func getRequests(c echo.Context) error {
	rows, err := db.Query(`SELECT id, content FROM requests`)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to fetch requests from database")
	}
	defer rows.Close()

	var requests []Request
	for rows.Next() {
		var request Request
		err := rows.Scan(&request.ID, &request.Content)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to fetch request from database")
		}
		requests = append(requests, request)
	}

	return c.JSON(http.StatusOK, requests)
}

func getRequestByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID")
	}

	var request Request
	err = db.QueryRow(`SELECT id, content FROM requests WHERE id = $1`, id).Scan(&request.ID, &request.Content)
	if err == sql.ErrNoRows {
		return c.String(http.StatusNotFound, "Request not found")
	} else if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to fetch request from database")
	}

	return c.JSON(http.StatusOK, request)
}


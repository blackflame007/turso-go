package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

// User represents a user in the leaderboard
type User struct {
	Name      string
	Email     string
	HighScore int
}

func main() {
	// Setup dot env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var dbUrl = fmt.Sprintf("%s?authToken=%s", os.Getenv("DB_URL"), os.Getenv("DB_AUTH_TOKEN"))
	db, err := sql.Open("libsql", dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", dbUrl, err)
		os.Exit(1)
	}
	defer db.Close()

	// Create users table
	createTable(db)

	// Insert a user
	insertUser(db, "John Doe", "john@example.com", 100)
	insertUser(db, "Bobs Burgers", "bob@example.com", 70)
	insertUser(db, "Jane Doe", "jane@example.com", 90)

	// Retrieve and display the leaderboard
	users, err := getLeaderboard(db)
	if err != nil {
		log.Fatalf("Failed to get leaderboard: %s", err)
	}
	displayLeaderboard(users)
}

func createTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		high_score INT NOT NULL
	);`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %s", err)
	}
}

func insertUser(db *sql.DB, name string, email string, highScore int) {
	// Check if the user already exists
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", email).Scan(&exists)
	if err != nil {
		log.Fatalf("Failed to check if user exists: %s", err)
	}

	if exists {
		fmt.Println("User with this email already exists")
		return
	}

	// Insert the user
	insertUserSQL := `INSERT INTO users (name, email, high_score) VALUES (?, ?, ?);`
	_, err = db.Exec(insertUserSQL, name, email, highScore)
	if err != nil {
		log.Fatalf("Failed to insert user: %s", err)
	}
}

func getLeaderboard(db *sql.DB) ([]User, error) {
	query := `SELECT name, email, high_score FROM users ORDER BY high_score DESC;`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.Name, &u.Email, &u.HighScore); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func displayLeaderboard(users []User) {
	fmt.Println("Leaderboard:")
	for i, user := range users {
		fmt.Printf("%d. %s (Email: %s) - Score: %d\n", i+1, user.Name, user.Email, user.HighScore)
	}
}

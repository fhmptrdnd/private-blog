package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Path ke database di cmd/web folder
	db, err := sql.Open("sqlite3", "../web/blog.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Query semua artikel
	rows, err := db.Query(`
		SELECT id, title, author, created_at, updated_at, views, owner_id, deleted_at
		FROM articles
		ORDER BY created_at DESC
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("=== DATABASE CONTENTS ===\n")
	count := 0
	for rows.Next() {
		var id, title, author, ownerID string
		var createdAt, updatedAt string
		var views int
		var deletedAt sql.NullString

		err := rows.Scan(&id, &title, &author, &createdAt, &updatedAt, &views, &ownerID, &deletedAt)
		if err != nil {
			log.Fatal(err)
		}

		count++
		fmt.Printf("--- Article #%d ---\n", count)
		fmt.Printf("ID:         %s\n", id)
		fmt.Printf("Title:      %s\n", title)
		fmt.Printf("Author:     %s\n", author)
		fmt.Printf("Created:    %s\n", createdAt)
		fmt.Printf("Updated:    %s\n", updatedAt)
		fmt.Printf("Views:      %d\n", views)
		fmt.Printf("Owner:      %s\n", ownerID)
		if deletedAt.Valid {
			fmt.Printf("Deleted:    %s [SOFT DELETED]\n", deletedAt.String)
		} else {
			fmt.Printf("Deleted:    [ACTIVE]\n")
		}
		fmt.Println()
	}

	fmt.Printf("Total: %d articles\n", count)
}

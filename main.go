package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/lib/pq"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("Error: DATABASE_URL is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Println("Database connection warning:", err)
	} else {
		fmt.Println("Database Connected Successfully!")
	}

	app := fiber.New()
	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("ONLYMAX Server is Running!")
	})

	app.Get("/users", func(c *fiber.Ctx) error {
		rows, err := db.Query("SELECT id, username, email FROM profiles")
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		defer rows.Close()

		var users []User = []User{} // กำหนดเป็น slice ว่างเพื่อไม่ให้ส่งค่า null กลับไป
		for rows.Next() {
			var u User
			if err := rows.Scan(&u.ID, &u.Username, &u.Email); err != nil {
				return err
			}
			users = append(users, u)
		}
		return c.JSON(users)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "10000"
	}
	
	fmt.Printf("Server is starting on port %s...\n", port)
	log.Fatal(app.Listen(":" + port))
}

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

// 1. ย้าย Struct มาไว้นอก main (ถูกต้องแล้ว)
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

	// แก้ให้เป็น Warning เพื่อไม่ให้แอปตายตอน Build
	if err := db.Ping(); err != nil {
		log.Println("Database connection warning:", err)
	}

	app := fiber.New()
	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("ONLYMAX Server is Running!")
	})

	// --- ตรวจสอบตรงนี้: ต้องอยู่ภายในปีกกาของ main เท่านั้น ---
	app.Get("/users", func(c *fiber.Ctx) error {
		rows, err := db.Query("SELECT id, username, email FROM profiles")
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		defer rows.Close()

		var users []User
		for rows.Next() {
			var u User
			if err := rows.Scan(&u.ID, &u.Username, &u.Email); err != nil {
				return err
			}
			users = append(users, u)
		}
		return c.JSON(users)
	})
	// --------------------------------------------------

	port := os.Getenv("PORT")
	if port == "" {
		port = "10000"
	}
	log.Fatal(app.Listen(":" + port))
} // ปิดท้ายฟังก์ชัน main (ตรวจเช็คปีกกาอันนี้ให้ดีครับ)

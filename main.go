package main


import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/lib/pq" // Driver สำหรับ PostgreSQL
)

func main() {
	// 1. อ่านค่า Connection String จาก Environment Variable ของเครื่อง Server
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("Error: DATABASE_URL environment variable is not set")
	}

	// 2. เชื่อมต่อ Database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// ทดสอบว่าเชื่อมต่อได้จริงไหม
	if err := db.Ping(); err != nil {
		log.Fatal("Database unreachable:", err)
	}
	fmt.Println("Database Connected Successfully!")

	app := fiber.New()
	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("ONLYMAX Server is Running with Database!")
	})

	// ... ใส่ Code API อื่นๆ ต่อตรงนี้ ...

	// Render จะกำหนด PORT มาให้ เราต้องอ่านค่า Port นั้น
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(":" + port))
}

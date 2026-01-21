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


// โครงสร้างข้อมูลสำหรับรับ/ส่ง JSON
type User struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}

// ... ในส่วน main function หลัง app := fiber.New() ...

// 1. ดึงข้อมูล User ทั้งหมด (GET)
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

// 2. เพิ่ม User ใหม่ (POST)
app.Post("/users", func(c *fiber.Ctx) error {
    u := new(User)
    if err := c.BodyParser(u); err != nil {
        return c.Status(400).SendString("Invalid input")
    }

    query := "INSERT INTO profiles (username, email) VALUES ($1, $2) RETURNING id"
    err := db.QueryRow(query, u.Username, u.Email).Scan(&u.ID)
    if err != nil {
        return c.Status(500).SendString(err.Error())
    }

    return c.Status(201).JSON(u)
})

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
    log.Println("Database connection warning:", err) // แจ้งเตือนใน Log แต่ไม่ปิดแอป
	} else {
    fmt.Println("Database Connected Successfully!")
	}

	app := fiber.New()
	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("ONLYMAX Server is Running with Database!")
	})

	// ... ใส่ Code API อื่นๆ ต่อตรงนี้ ...

	// Render จะกำหนด PORT มาให้ เราต้องอ่านค่า Port นั้น
	port := os.Getenv("PORT")
	if port == "" {
    port = "10000" // ค่าเริ่มต้นถ้าไม่มี env
	}
	app.Listen(":" + port)
}

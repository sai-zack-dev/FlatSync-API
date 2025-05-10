package handlers

import (
	"database/sql"
	"time"

	"my-fiber-api/database"
	"my-fiber-api/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("your_secret_key") // Use os.Getenv("JWT_SECRET") with godotenv if using .env

func Register(c *fiber.Ctx) error {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var body Request

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(body.Password), 14)

	stmt, err := database.DB.Prepare("INSERT INTO users(email, password) VALUES (?, ?)")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "DB error"})
	}
	_, err = stmt.Exec(body.Email, hash)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Email already exists"})
	}

	return c.JSON(fiber.Map{"message": "Registered successfully"})
}

func Login(c *fiber.Ctx) error {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var body Request
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	var user models.User
	row := database.DB.QueryRow("SELECT id, password FROM users WHERE email = ?", body.Email)
	err := row.Scan(&user.ID, &user.Password)
	if err == sql.ErrNoRows {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)) != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Incorrect password"})
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(jwtSecret)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not login"})
	}

	return c.JSON(fiber.Map{"token": signed})
}

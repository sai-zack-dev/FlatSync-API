package handlers

import (
	"database/sql"
	"time"

	"github.com/sai-zack-dev/FlatSync-API/database"
	"github.com/sai-zack-dev/FlatSync-API/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"

	"os"

	"github.com/joho/godotenv"
)

var jwtSecret []byte

func init() {
	godotenv.Load()
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
}

func Register(c *fiber.Ctx) error {
	type Request struct {
		Name     string `json:"name"`     // Required field
		Email    string `json:"email"`    // Required field
		Password string `json:"password"` // Required field
		Dob      string `json:"dob"`      // Nullable field
		Avatar   string `json:"avatar"`   // Nullable field
	}

	var body Request

	// Parse the request body
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	// Validate that the required fields are present
	if body.Name == "" || body.Email == "" || body.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Name, email, and password are required"})
	}

	// Hash the password
	hash, _ := bcrypt.GenerateFromPassword([]byte(body.Password), 14)

	// Insert the user into the database
	stmt, err := database.DB.Prepare(`
        INSERT INTO users(name, email, password, dob, avatar) 
        VALUES (?, ?, ?, ?, ?)
    `)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "DB error"})
	}
	_, err = stmt.Exec(body.Name, body.Email, hash, body.Dob, body.Avatar)
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

	// Parse the request body
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	// Add validation for empty email or password
	if body.Email == "" || body.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Email and password required"})
	}

	// Retrieve the user from the database
	var user models.User
	row := database.DB.QueryRow("SELECT id, password FROM users WHERE email = ?", body.Email)
	err := row.Scan(&user.ID, &user.Password)
	if err == sql.ErrNoRows {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	// Compare the provided password with the hashed password in the database
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)) != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Incorrect password"})
	}

	// Generate JWT token
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(jwtSecret)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not login"})
	}

	return c.JSON(fiber.Map{"message": "Login successfully", "token": signed})
}

func Protected(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or invalid token"})
	}
	tokenString := authHeader[7:]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
	}

	userID := claims["user_id"]
	return c.JSON(fiber.Map{
		"message": "Access granted to protected route",
		"user_id": userID,
	})
}

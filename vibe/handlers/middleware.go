package handlers

import (
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func (h *Handler) AuthMiddleware(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")

	// Remove the "Bearer " prefix from the token string
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(h.JWTSigningKey), nil
	})

	if err != nil || !token.Valid {
		logrus.Errorf(err.Error())
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	// Get the user ID from the token's claims
	var userID float64
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if id, exists := claims["user_id"]; exists {
			userID = id.(float64)
			c.Locals("user_id", userID) // Set the user ID in locals for other handlers to access
		}
		if email, exists := claims["email"]; exists {
			c.Locals("user_email", email) // Set the user ID in locals for other handlers to access
		}
		if name, exists := claims["name"]; exists {
			c.Locals("user_name", name) // Set the user ID in locals for other handlers to access
		}
		if organizationID, exists := claims["organization_id"]; exists {
			c.Locals("organization_id", organizationID) // Set the user ID in locals for other handlers to access
		}

	}

	// Check if the token is still valid (not invalidated by logout)

	// Token is valid, proceed to the next handler
	return c.Next()
}

// TODOAuthMiddleware is SHOULD never be used, it bypasses the auth middleware
// ONLY use this for testing purposes
func (h *Handler) TODOAuthMiddleware(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")

	// Remove the "Bearer " prefix from the token string
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Parse the token
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(h.JWTSigningKey), nil
	})

	// if  err != nil || !token.Valid {
	// 	logrus.Errorf(err.Error())
	// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
	// 		"message": "Unauthorized",
	// 	})
	// }

	// Get the user ID from the token's claims
	var userID float64
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if id, exists := claims["user_id"]; exists {
			userID = id.(float64)
			c.Locals("user_id", userID) // Set the user ID in locals for other handlers to access
		}
		if email, exists := claims["email"]; exists {
			c.Locals("user_email", email) // Set the user ID in locals for other handlers to access
		}
		if name, exists := claims["name"]; exists {
			c.Locals("user_name", name) // Set the user ID in locals for other handlers to access
		}
		if organizationID, exists := claims["organization_id"]; exists {
			c.Locals("organization_id", organizationID) // Set the user ID in locals for other handlers to access
		}

	}

	// Check if the token is still valid (not invalidated by logout)

	// Token is valid, proceed to the next handler
	return c.Next()
}

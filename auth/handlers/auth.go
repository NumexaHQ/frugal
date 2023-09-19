package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/NumexaHQ/captainCache/model"
	postgresql_db "github.com/NumexaHQ/captainCache/numexa-common/postgresql/postgresql-db"
	"github.com/NumexaHQ/captainCache/pkg/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const keyLength = 32

// generateAPIKey generates a random string of length 32 with prefix "sk"
func generateAPIKey() string {
	b := make([]byte, keyLength)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return fmt.Sprintf("sk-%s", string(b))
}

// func hashPassword(password string) string {
// 	hashPassword := utils.HashPassword(password)
// }

func (h *Handler) CreateApiKey(c *fiber.Ctx) error {
	type RequestBody postgresql_db.NxaApiKey
	var reqBody RequestBody
	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if reqBody.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Name cannot be empty",
		})
	}

	userID := c.Locals("user_id").(float64)
	_, err := h.DB.GetUserById(c.Context(), int32(userID))
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Errorf("error getting user by id: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Invalid userId",
			})
		}
	}

	apiKey := generateAPIKey()
	//hashedApiKey := hashAPIKey(apiKey)

	// todo: get user_id from jwt token
	// todo: get project_id from request body
	// validate if user is part of project
	userEmail := c.Locals("user_email").(string)
	user, err := h.DB.GetUserByEmail(c.Context(), userEmail)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Errorf("error getting user by email: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error getting user by email",
			})
		} else {
			log.Errorf("user not found: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Something went wrong. Please contact the administrator",
			})
		}
	}

	_, err = h.DB.GetProjectUserByProjectIDAndUserID(c.Context(), reqBody.ProjectID, user.ID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Errorf("error getting project user by projectID and userID: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error getting project user by projectID and userID",
			})
		} else {
			log.Errorf("user is not part of project: %v", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "User is not part of project",
			})
		}
	}

	// check if api key with that name exists within the project
	k, err := h.DB.GetAPIKeyByNameAndProjectId(c.Context(), reqBody.Name, reqBody.ProjectID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Errorf("error getting api key by name and projectID: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error getting api key by name and projectID",
			})
		}
	}

	if k.Name != "" {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "API key with that name already exists",
		})
	}

	_, err = h.DB.CreateApiKey(c.Context(), postgresql_db.NxaApiKey{
		Name:      reqBody.Name,
		ApiKey:    apiKey,
		UserID:    user.ID,
		ProjectID: reqBody.ProjectID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 365),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		log.Errorf("error generating api key: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error generated api key",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "API key generated successfully. Please store it securely!",
		"key":     apiKey,
	})
}

func (h *Handler) RegisterHandler(c *fiber.Ctx) error {
	var organization postgresql_db.Organization
	var project postgresql_db.Project
	var userReq model.RegisterUserRequest

	if err := c.BodyParser(&userReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if err := h.Validator.Struct(userReq); err != nil {
		log.Debugf("error validating user: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	user := postgresql_db.User{
		Name:     userReq.Name,
		Email:    userReq.Email,
		Password: userReq.Password,
	}

	u, err := h.DB.GetUserByEmail(c.Context(), user.Email)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Errorf("error getting user by email: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error getting user by email",
			})
		}
	}
	if u.Email != "" {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Username already exists",
		})
	}

	// create organization
	organization.Name = utils.GenerateOrganizationName()

	log.Infof("organization: %+v", organization)
	organization, err = h.DB.CreateOrganization(c.Context(), organization)
	if err != nil {
		log.Errorf("error creating organization: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error creating organization",
		})
	}

	// create project
	project.Name = utils.GenerateProjectName()
	project.OrganizationID = organization.ID

	project, err = h.DB.CreateProject(c.Context(), project)
	if err != nil {
		log.Errorf("error creating project: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error creating project",
		})
	}

	user.OrganizationID = organization.ID

	// Hashing the password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Errorf("error hashing password: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error hashing password",
		})
	}

	user.Password = string(passwordHash)

	// create user
	user, err = h.DB.RegisterUser(c.Context(), user)
	if err != nil {
		// todo
		log.Errorf("error registering user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error registering user. Please contact administrator",
		})
	}

	// create project user
	projectUser := postgresql_db.ProjectUser{
		ProjectID: project.ID,
		UserID:    user.ID,
		RoleID:    1,
	}

	_, err = h.DB.CreateProjectUser(c.Context(), projectUser)
	if err != nil {
		log.Errorf("error creating project user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error creating project user",
		})
	}

	tokenString, err := generateJWTToken(user, h.JWTSigningKey)

	if err != nil {
		log.Errorf("error generating token: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error generating token",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":    "User registered successfully",
		"token":      tokenString,
		"project_id": project.ID,
	})
}

func (h *Handler) LoginHandler(c *fiber.Ctx) error {
	var user postgresql_db.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// Get the user from Postgres
	u, err := h.DB.GetUserByEmail(c.Context(), user.Email)
	if err != nil {
		log.Errorf("error getting user by email: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid username or password",
		})
	}

	// Check if the password is correct
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password)); err != nil {
		log.Errorf("error comparing password: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid username or password",
		})
	}

	project, err := h.DB.GetProjectsByOrgId(c.Context(), u.OrganizationID)
	if err != nil {
		log.Errorf("error getting projects by org id: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error getting projects by org id",
		})
	}

	// Generate a JWT token
	tokenString, err := generateJWTToken(u, h.JWTSigningKey)

	if err != nil {
		log.Errorf("error generating token: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error generating token",
		})
	}

	// update last login
	err = h.DB.UpdateUserLastLogin(c.Context(), u)
	if err != nil {
		log.Errorf("error updating user last login: %v", err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token":      tokenString,
		"project_id": project[0].ID,
	})

}

func generateJWTToken(user postgresql_db.User, jwtSigningKey string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = user.Email
	claims["user_id"] = user.ID
	claims["name"] = user.Name
	claims["organization_id"] = user.OrganizationID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token valid for 24 hours

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(jwtSigningKey))

	return tokenString, err
}

// SHOULD ONLY BE USED FOR TESTING
// DONOT USE IN PRODUCTION
func (h *Handler) DummyAuthMiddleware(c *fiber.Ctx) error {
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

	// if err != nil {
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
			c.Locals("name", name) // Set the user ID in locals for other handlers to access
		}
		if organizationID, exists := claims["organization_id"]; exists {
			c.Locals("organization_id", organizationID) // Set the user ID in locals for other handlers to access
		}

	}

	// Check if the token is still valid (not invalidated by logout)

	// Token is valid, proceed to the next handler
	return c.Next()
}

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
			c.Locals("name", name) // Set the user ID in locals for other handlers to access
		}
		if organizationID, exists := claims["organization_id"]; exists {
			c.Locals("organization_id", organizationID) // Set the user ID in locals for other handlers to access
		}

	}

	// Check if the token is still valid (not invalidated by logout)

	// Token is valid, proceed to the next handler
	return c.Next()
}

func (h *Handler) LogoutHandler(c *fiber.Ctx) error {
	_ = c.Locals("user_email").(string)

	// Get the token from the Authorization header
	tokenString := c.Get("Authorization")
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// :todo

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}

// func (h *Handler) GoogleAuthCallback(c *fiber.Ctx) error {
// 	type Accesstoken struct {
// 		AccessToken string `json:"access_token"`
// 	}
// 	var token Accesstoken
// 	if err := c.BodyParser(&token); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"message": "Invalid request body",
// 		})
// 	}
// 	responseBody, err := http.Post("https://oauth2.googleapis.com/tokeninfo?access_token="+token.AccessToken, "application/json", nil)
// 	if err != nil {
// 		log.Errorf("error getting token info: %v", err)
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "Error during Google authentication",
// 		})
// 	}
// 	defer responseBody.Body.Close()

// 	var userInfo map[string]interface{}
// 	if err := json.NewDecoder(responseBody.Body).Decode(&userInfo); err != nil {
// 		log.Errorf("error decoding token info: %v", err)
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "Error during Google authentication",
// 		})
// 	}

// 	// Check if the user exists in the database
// 	existingUser, err := h.DB.GetUserByEmail(c.Context(), userInfo["email"].(string))
// 	if err != nil {
// 		if !errors.Is(err, sql.ErrNoRows) {
// 			log.Errorf("error getting user by email: %v", err)
// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 				"message": "Error getting user by email",
// 			})
// 		}
// 	}

// 	// If the user already exists, generate a JWT token and log them in
// 	if existingUser.Email != "" {
// 		tokenString, err := generateJWTToken(existingUser)
// 		if err != nil {
// 			log.Errorf("error generating token: %v", err)
// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 				"message": "Error generating token",
// 			})
// 		}

// 		return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 			"message": "User logged in successfully with Google authentication",
// 			"token":   tokenString,
// 		})
// 	}

// 	// User does not exist, create a new user
// 	var user postgresql_db.User
// 	var organization postgresql_db.Organization
// 	var project postgresql_db.Project

// 	// create organization
// 	organization.Name = utils.GenerateOrganizationName()

// 	log.Infof("organization: %+v", organization)
// 	organization, err = h.DB.CreateOrganization(c.Context(), organization)
// 	if err != nil {
// 		log.Errorf("error creating organization: %v", err)
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "Error creating organization",
// 		})
// 	}

// 	// create project
// 	project.Name = utils.GenerateProjectName()
// 	project.OrganizationID = organization.ID

// 	log.Infof("project: %+v", project)
// 	project, err = h.DB.CreateProject(c.Context(), project)
// 	if err != nil {
// 		log.Errorf("error creating project: %v", err)
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "Error creating project",
// 		})
// 	}

// 	user.OrganizationID = organization.ID
// 	user.Email = userInfo["email"].(string)
// 	user.Name = userInfo["name"].(string)
// 	user.Password = userInfo["id"].(string)
// 	log.Infof("user: %+v", user)

// 	// create user
// 	user, err = h.DB.RegisterUser(c.Context(), user)
// 	if err != nil {
// 		// todo
// 		log.Errorf("error registering user: %v", err)
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "Error registering user. Please contact administrator",
// 		})
// 	}

// 	// create project user
// 	projectUser := postgresql_db.ProjectUser{
// 		ProjectID: project.ID,
// 		UserID:    user.ID,
// 		RoleID:    1,
// 	}

// 	log.Infof("projectUser: %+v", projectUser)

// 	_, err = h.DB.CreateProjectUser(c.Context(), projectUser)
// 	if err != nil {
// 		log.Errorf("error creating project user: %v", err)
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "Error creating project user",
// 		})
// 	}

// 	tokenString, err := generateJWTToken(userInfo["email"].(string))
// 	if err != nil {
// 		log.Errorf("error generating token: %v", err)
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "Error generating token",
// 		})
// 	}

// 	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
// 		"message": "User registered successfully with Google authentication",
// 		"token":   tokenString,
// 	})
// }

func (h *Handler) GetAPIkeyByUserId(c *fiber.Ctx) error {
	userEmail := c.Locals("user_email").(string)
	user, err := h.DB.GetUserByEmail(c.Context(), userEmail)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Errorf("error getting user by email: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error getting user by email",
			})
		} else {
			log.Errorf("user not found: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Something went wrong. Please contact the administrator",
			})
		}
	}

	keys, err := h.DB.GetAllApiKeysByUserId(c.Context(), user.ID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Errorf("error getting api keys by user id: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error getting api keys by user id",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "API keys fetched successfully",
		"keys":    keys,
	})
}

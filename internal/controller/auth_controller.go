package controller

import (
	"github.com/budhilaw/personal-website-backend/config"
	"github.com/budhilaw/personal-website-backend/internal/model"
	"github.com/budhilaw/personal-website-backend/internal/service"
	"github.com/gofiber/fiber/v2"
)

// AuthController handles authentication-related requests
type AuthController struct {
	authService service.AuthService
	cfg         config.Config
}

// NewAuthController creates a new AuthController
func NewAuthController(authService service.AuthService, cfg config.Config) *AuthController {
	return &AuthController{
		authService: authService,
		cfg:         cfg,
	}
}

// Login handles login requests
func (c *AuthController) Login(ctx *fiber.Ctx) error {
	var loginReq model.UserLogin

	if err := ctx.BodyParser(&loginReq); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if loginReq.Username == "" || loginReq.Password == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username and password are required",
		})
	}

	// Login - pass the Fiber context for IP and user agent tracking
	resp, err := c.authService.Login(ctx.Context(), loginReq.Username, loginReq.Password, ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	return ctx.JSON(resp)
}

// GetProfile handles get profile requests
func (c *AuthController) GetProfile(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(string)

	profile, err := c.authService.GetProfile(ctx.Context(), userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get profile",
		})
	}

	return ctx.JSON(profile)
}

// UpdateProfile handles update profile requests
func (c *AuthController) UpdateProfile(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(string)

	var profileReq model.ProfileUpdate
	if err := ctx.BodyParser(&profileReq); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if profileReq.FirstName == "" || profileReq.Email == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "First name and email are required",
		})
	}

	if err := c.authService.UpdateProfile(ctx.Context(), userID, &profileReq); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update profile",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Profile updated successfully",
	})
}

// UpdateAvatar handles update avatar requests
func (c *AuthController) UpdateAvatar(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(string)

	// Get avatar from form field
	avatar := ctx.FormValue("avatar")
	if avatar == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Avatar is required",
		})
	}

	if err := c.authService.UpdateAvatar(ctx.Context(), userID, avatar); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update avatar",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Avatar updated successfully",
	})
}

// UpdatePassword handles update password requests
func (c *AuthController) UpdatePassword(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(string)

	// Parse request
	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if req.CurrentPassword == "" || req.NewPassword == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Current password and new password are required",
		})
	}

	if err := c.authService.UpdatePassword(ctx.Context(), userID, req.CurrentPassword, req.NewPassword); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password updated successfully",
	})
}

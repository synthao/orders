package sso

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/synthao/orders/internal/module/sso"
	"strings"
)

var ErrExtractToken = fmt.Errorf("failed to extract token")

type Config struct {
	Client *sso.Client
}

func New(config Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		token, err := extractToken(token)
		if err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		isAuthorized, err := config.Client.IsAuthorized(token)
		if err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		if !isAuthorized {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		return c.Next()
	}
}

func extractToken(s string) (string, error) {
	chunks := strings.Split(s, " ")
	if len(chunks) < 2 {
		return "", ErrExtractToken
	}

	return chunks[1], nil
}

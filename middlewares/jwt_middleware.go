package middlewares

import (
	"log"
	"os"
	"strings"

	"hrms_go/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTProtected() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		auth := ctx.Get("Authorization")
		if auth == "" {
			return utils.Error(ctx, 401, "missing token")
		}

		tokenStr := strings.Replace(auth, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			return utils.Error(ctx, 401, "invalid token")
		}

		claims := token.Claims.(jwt.MapClaims)
		log.Println("JWT Claims:", claims)
		
		ctx.Locals("employee_nik", claims["employee_nik"])
		ctx.Locals("user_id", claims["user_id"])

		return ctx.Next()
	}
}

package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/okm321/mahking-go/pkg/logger"
)

func Recover() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				logger.ErrorContext(c.Context(), fmt.Sprintf("%v", r))

				// エラーに変換
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}

				// Fiberのエラーハンドラーに委譲
				// （直接レスポンスを返さない）
				_ = c.App().Config().ErrorHandler(c, err)
			}
		}()

		return c.Next()
	}
}

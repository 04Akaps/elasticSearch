package v1

import "github.com/gofiber/fiber/v2"

type v1 struct {
	// TODO service
}

func RegisterV1Router(r fiber.Router) {
	v := v1{}

	r.Get("/read-test", v.readTest)

}

func (v v1) readTest(c *fiber.Ctx) error {

	return nil
}

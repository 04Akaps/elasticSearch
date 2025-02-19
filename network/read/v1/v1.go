package v1

import (
	"github.com/04Akaps/elasticSearch.git/common/validator"
	"github.com/04Akaps/elasticSearch.git/types/request"
	"github.com/gofiber/fiber/v2"

	v1Service "github.com/04Akaps/elasticSearch.git/service/v1"
)

type v1 struct {
	service v1Service.V1
}

func RegisterV1Router(r fiber.Router, service v1Service.V1) {
	v := v1{service}

	r.Get("/read-test", v.readTest)

}

func (v v1) readTest(c *fiber.Ctx) error {
	var req request.ReadTestRequest

	err := validator.QueryParser(&req, c)

	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(err)
	}

	v.service.ReadTest(req)

	return nil
}

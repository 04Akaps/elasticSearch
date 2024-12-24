package v1

import (
	"github.com/04Akaps/elasticSearch.git/common/validator"
	v1Service "github.com/04Akaps/elasticSearch.git/service/v1"
	"github.com/04Akaps/elasticSearch.git/types/request"
	"github.com/gofiber/fiber/v2"
)

type v1 struct {
	service v1Service.V1
}

func RegisterV1Router(r fiber.Router, service v1Service.V1) {
	v := v1{service}

	r.Post("/insert-test", v.insertTest)
	r.Post("/insert-mapper", v.insertMapper)
}

func (v v1) insertTest(c *fiber.Ctx) error {
	var req request.InsertTestRequest

	err := validator.BodyParser(&req, c)

	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(err)
	}

	v.service.InsertTest(req)

	return nil
}

func (v v1) insertMapper(c *fiber.Ctx) error {
	v.service.InsertMapperTest()
	return nil
}

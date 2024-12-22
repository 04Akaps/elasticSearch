package insert

import (
	v1 "github.com/04Akaps/elasticSearch.git/network/insert/v1"
	"github.com/04Akaps/elasticSearch.git/service"
	"github.com/gofiber/fiber/v2"
)

func RegisterInsertRouter(r fiber.Router, service service.Manager) {
	v1.RegisterV1Router(r.Group("/v1"), service.V1())
}

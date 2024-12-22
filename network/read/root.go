package read

import (
	v1 "github.com/04Akaps/elasticSearch.git/network/read/v1"
	"github.com/gofiber/fiber/v2"
)

func RegisterReadRouter(r fiber.Router) {

	v1.RegisterV1Router(r.Group("/v1"))
}

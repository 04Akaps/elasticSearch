package network

import (
	"context"
	"fmt"
	"github.com/04Akaps/elasticSearch.git/config"
	"github.com/04Akaps/elasticSearch.git/docker"
	insert "github.com/04Akaps/elasticSearch.git/network/Insert"
	"github.com/04Akaps/elasticSearch.git/network/read"
	"github.com/04Akaps/elasticSearch.git/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/fx"
	"log"
	"runtime/debug"
	"strings"
	"time"
)

type Router struct {
	engine *fiber.App
	port   string
}

var AllowHeaders = []string{
	"ORIGIN",
	"Content-Length",
	"Content-Type",
	"Access-Control-Allow-Headers",
	"Access-Control-Allow-Origin",
	"Authorization",
	"Cache-Control",
}

func NewRouter(
	lc fx.Lifecycle,
	config config.Config,
	manager service.Manager,
) Router {
	r := Router{
		port: fmt.Sprintf(":%s", config.Server.Port),
	}

	r.engine = fiber.New()
	r.engine.Use(
		recover.New(recover.Config{
			EnableStackTrace: true,
			StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
				msg := fmt.Sprintf("panic: %v\n%s\n", e, debug.Stack())
				log.Println(msg)
			},
		}),
		logger.New(logger.Config{
			Format:     "[FIBER] ${time} | ${status} | ${latency} | ${ip} | ${method} | ${cerr} | ${stack} | \"${path}\"\n",
			TimeFormat: "2006/01/02 - 15:04:05",
			TimeZone:   "Local",
		}),
		cors.New(cors.Config{
			AllowOrigins:     "*",
			AllowMethods:     strings.Join([]string{"GET", "POST", "PUT", "DELETE", "PATCH"}, ", "),
			AllowHeaders:     strings.Join(AllowHeaders, ", "),
			ExposeHeaders:    strings.Join(AllowHeaders, ", "),
			AllowCredentials: false,
			AllowOriginsFunc: func(origin string) bool { return true },
			MaxAge:           12 * int(time.Hour.Seconds()),
		}),
	)

	r.engine.Get("/healthCheck", r.healthCheck)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			docker.Initialize(config)
			go func() {
				log.Println("server start", "endpoint", r.port)
				if err := r.engine.Listen(r.port); err != nil {
					log.Println("Error Starting Server", "cerr", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Server Shutdown", "cerr", ctx.Err())
			return r.engine.Shutdown()
		},
	})

	read.RegisterReadRouter(r.group("/api/read"), manager)
	insert.RegisterInsertRouter(r.group("/api/insert"), manager)

	return r
}

func (r Router) group(path string) fiber.Router {
	return r.engine.Group(path)
}

func (r Router) healthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ok",
	})
}

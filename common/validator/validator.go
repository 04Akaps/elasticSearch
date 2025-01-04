package validator

import (
	"fmt"
	"github.com/04Akaps/elasticSearch.git/types/protocol"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"log"
)

const (
	_validationError = "validation cerr: %s"
)

type RequestValidator struct {
	validator *validator.Validate
}

var GlobalValidator RequestValidator

func init() {
	GlobalValidator = RequestValidator{
		validator: validator.New(),
	}

}

func RegisterValidation(tag string, fn func(fl validator.FieldLevel) bool) error {
	return GlobalValidator.validator.RegisterValidation(tag, fn)
}

func (r *RequestValidator) Validate(data interface{}) []ErrValidation {
	var validationErrors []ErrValidation
	errs := r.validator.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			fmt.Println(err.Field(), err.Tag(), err.Value(), err.Namespace())
			v := NewErrValidation(err.Field(), err.Tag(), err.Value())
			validationErrors = append(validationErrors, *v)
		}
	}

	return validationErrors
}

func QueryParser(req interface{}, c *fiber.Ctx) error {
	if err := c.QueryParser(req); err != nil {
		// TODO Log, custom cerr
		return err
	}

	errs := GlobalValidator.Validate(req)

	if len(errs) > 0 {
		msg := fmt.Sprintf(_validationError, errs[0].Error())
		log.Println("Failed to parsing query request")
		return protocol.Response(msg, protocol.FailedQueryParsing)
	}

	return nil
}

func BodyParser(req interface{}, c *fiber.Ctx) error {
	if err := c.BodyParser(req); err != nil {
		// TODO Log, custom cerr
		return nil
	}

	errs := GlobalValidator.Validate(req)

	if len(errs) > 0 {
		msg := fmt.Sprintf(_validationError, errs[0].Error())
		log.Println("Failed to parsing body")
		return protocol.Response(msg, protocol.FailedBodyParsing)
	}

	return nil
}

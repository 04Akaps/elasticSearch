package json

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2/log"
	"reflect"
)

type handler struct {
	marshal   func(val interface{}) ([]byte, error)
	unmarshal func(buf []byte, val interface{}) error
}

var JsonHandler handler

func init() {
	JsonHandler = handler{
		marshal:   sonic.Marshal,
		unmarshal: sonic.Unmarshal,
	}
}

func (h handler) Marshal(v interface{}) ([]byte, error) {
	bytes, err := h.marshal(v)

	if err != nil {
		log.Errorf("Failed to marshal value", "type", reflect.TypeOf(v).String(), "err", err)
		return nil, err
	}

	return bytes, nil
}

func (h handler) Unmarshal(buffer []byte, v interface{}) error {
	err := h.unmarshal(buffer, v)

	if err != nil {
		log.Errorf("Failed to unmarshal value", "buffer", string(buffer), "type", reflect.TypeOf(v).String(), "err", err)
		return err
	}

	return nil
}

func (h handler) Handle(buf interface{}, v interface{}) error {

	bytes, err := h.marshal(buf)

	if err != nil {
		return err
	}

	err = h.unmarshal(bytes, v)

	if err != nil {
		return err
	}

	return nil
}

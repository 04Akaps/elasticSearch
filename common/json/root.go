package json

import (
	"bytes"
	"github.com/04Akaps/elasticSearch.git/common/sync"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2/log"
	"reflect"
)

type handler struct {
	marshal    func(val interface{}) ([]byte, error)
	unmarshal  func(buf []byte, val interface{}) error
	memoryPool *sync.Pool[*bytes.Buffer]
}

var JsonHandler handler

func init() {
	JsonHandler = handler{
		marshal:   sonic.Marshal,
		unmarshal: sonic.Unmarshal,
		memoryPool: sync.NewPool[*bytes.Buffer](func() *bytes.Buffer {
			return new(bytes.Buffer)
		}),
	}
}

func (h handler) Marshal(v interface{}) ([]byte, error) {
	memoryBuf := h.memoryPool.Get()
	defer h.memoryPool.Put(memoryBuf)

	memoryBuf.Reset()
	data, err := h.marshal(v)
	memoryBuf.Write(data)

	if err != nil {
		log.Errorf("Failed to marshal value", "type", reflect.TypeOf(v).String(), "cerr", err)
		return nil, err
	}

	return memoryBuf.Bytes(), nil
}

func (h handler) Unmarshal(buffer []byte, v interface{}) error {
	memoryBuf := h.memoryPool.Get()
	defer h.memoryPool.Put(memoryBuf)

	memoryBuf.Reset()
	memoryBuf.Write(buffer)

	err := h.unmarshal(memoryBuf.Bytes(), v)

	if err != nil {
		log.Errorf("Failed to unmarshal value", "buffer", string(buffer), "type", reflect.TypeOf(v).String(), "cerr", err)
		return err
	}

	return nil

}

func (h handler) Handle(buf interface{}, v interface{}) error {

	data, err := h.marshal(buf)

	if err != nil {
		return err
	}

	err = h.unmarshal(data, v)

	if err != nil {
		return err
	}

	return nil
}

package request

type InsertTestRequest struct {
	Index string `json:"index" validate:"required"`
	Key   string `json:"key" validate:"required"`
	Value string `json:"value" validate:"required"`
}

type ReadTestRequest struct {
	Index string `query:"index" validate:"required"`
	Key   string `query:"key" validate:"required"`
	Value string `query:"value" validate:"required"`
}

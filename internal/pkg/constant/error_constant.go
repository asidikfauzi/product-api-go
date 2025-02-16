package constant

import "errors"

var (
	CategoryNotFound               = errors.New("category not found")
	InvalidJsonPayload             = errors.New("invalid request payload")
	InvalidQueryParameters         = errors.New("invalid query parameters")
	FailedToLoadTimeZone           = errors.New("failed to load timezone")
	UnexpectedError                = errors.New("unexpected error")
	JsonUnmarshalError             = errors.New("json unmarshal error")
	JsonMarshalError               = errors.New("json marshal error")
	KeyRedisNotExists              = errors.New("key not exists")
	RedisGetError                  = errors.New("redis get error")
	RedisSetError                  = errors.New("redis set error")
	RedisDeleteError               = errors.New("redis delete error")
	SomeCategoryNotFound           = errors.New("some categories not found")
	CategoryUnprocessableEntity    = errors.New("category unprocessable entity")
	CategoryAlreadyExists          = errors.New("category already exists")
	MeasurementNotFound            = errors.New("measurement not found")
	MeasurementUnprocessableEntity = errors.New("measurement unprocessable entity")
	MeasurementAlreadyExists       = errors.New("measurement already exists")
	ProductNotFound                = errors.New("product not found")
	ProductUnprocessableEntity     = errors.New("product unprocessable entity")
	ProductAlreadyExists           = errors.New("product already exists")
)

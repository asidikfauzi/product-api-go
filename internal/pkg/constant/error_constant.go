package constant

const (
	InvalidJsonPayload     string = "invalid request payload"
	InvalidQueryParameters string = "invalid query parameters"
	FailedToLoadTimeZone   string = "failed to load timezone: %w"

	CategoryNotFound            string = "category not found"
	SomeCategoryNotFound        string = "some categories not found"
	CategoryUnprocessableEntity string = "category unprocessable entity"
	CategoryAlreadyExists       string = "category already exists"

	MeasurementNotFound            string = "measurement not found"
	MeasurementUnprocessableEntity string = "measurement unprocessable entity"
	MeasurementAlreadyExists       string = "measurement already exists"

	ProductNotFound            string = "product not found"
	ProductUnprocessableEntity string = "product unprocessable entity"
	ProductAlreadyExists       string = "product already exists"
)

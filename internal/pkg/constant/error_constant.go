package constant

const (
	InvalidJsonPayload     string = "invalid request payload"
	InvalidQueryParameters string = "invalid query parameters"
	FailedToLoadTimeZone   string = "failed to load timezone: %w"

	CategoryNotFound            string = "category not found"
	CategoryUnprocessableEntity string = "category unprocessable entity"
	CategoryAlreadyExists       string = "category already exists"
)

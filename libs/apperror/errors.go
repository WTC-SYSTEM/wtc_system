package apperror

type ErrorsEnum = string

const (
	// System Error (Something went wrong...)
	WTC_000001 ErrorsEnum = "WTC-000001"
	// Bad Request (some thing wrong with user data)
	WTC_000002 ErrorsEnum = "WTC-000002"
	// Not Found (Not found)
	WTC_000003 ErrorsEnum = "WTC-000003"
	// User with such email is already registered
	WTC_000004 ErrorsEnum = "WTC-000004"
	// Invalid email or password
	WTC_000005 ErrorsEnum = "WTC-000005"

	WTC_000006 ErrorsEnum = "WTC-000006"
)

package constants

// Success Codes
const (
	SUCCESS_CODE    = "200"
	CREATED_CODE    = "201"
	ACCEPTED_CODE   = "202"
	NO_CONTENT_CODE = "204"
)

// Client Error Codes
const (
	BAD_REQUEST_CODE     = "400"
	UNAUTHORIZED_CODE    = "401"
	FORBIDDEN_CODE       = "403"
	NOT_FOUND_CODE       = "404"
	METHOD_NOT_ALLOWED   = "405"
	CONFLICT_CODE        = "409"
	UNPROCESSABLE_ENTITY = "422"
)

// Server Error Codes
const (
	INTERNAL_SERVER_ERROR = "500"
	NOT_IMPLEMENTED_CODE  = "501"
	BAD_GATEWAY_CODE      = "502"
	SERVICE_UNAVAILABLE   = "503"
	GATEWAY_TIMEOUT_CODE  = "504"
)

// Custom Application Codes
const (
	VALIDATION_ERROR_CODE = "1001"
	BUSINESS_ERROR_CODE   = "1002"
	DATABASE_ERROR_CODE   = "1003"
	EMPTY_CODE            = ""
	EMPTY_STRING          = ""
	EMPTY_INT             = 0
	TRUE_VALUE            = "true"
	FALSE_VALUE           = "false"
	EMPTY_VALUE           = ""
	USER                  = "user"

	// Role
	SuperAdminRole    = "Super Admin"
	CustomerRole      = "Customer"
	MerchantRole      = "Merchant"
	MerchantAdminRole = "Merchant Admin"
	AccessDenied      = "Access denied. You don't have permission"
)

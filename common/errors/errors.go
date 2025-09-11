package errors

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

type ErrorCode int

const (
	// General errors
	ErrCodeInternal ErrorCode = iota + 1000
	ErrCodeBadRequest
	ErrCodeUnauthorized
	ErrCodeForbidden
	ErrCodeNotFound
	ErrCodeConflict
	ErrCodeValidation
	ErrCodeTimeout
	ErrCodeRateLimit

	// Authentication errors
	ErrCodeInvalidCredentials ErrorCode = iota + 2000
	ErrCodeTokenExpired
	ErrCodeTokenInvalid
	ErrCodeUserNotFound
	ErrCodeUserAlreadyExists

	// Document errors
	ErrCodeDocumentNotFound ErrorCode = iota + 3000
	ErrCodeDocumentAccessDenied
	ErrCodeDocumentLocked
	ErrCodeInvalidDocumentFormat
	ErrCodeDocumentSizeExceeded

	// Workspace errors
	ErrCodeWorkspaceNotFound ErrorCode = iota + 4000
	ErrCodeWorkspaceAccessDenied
	ErrCodeInvitationExpired
	ErrCodeMaxMembersExceeded

	// Collaboration errors
	ErrCodeSessionNotFound ErrorCode = iota + 5000
	ErrCodeOperationConflict
	ErrCodeConcurrentEdit
)

type AppError struct {
	Code       ErrorCode      `json:"code"`
	Message    string         `json:"message"`
	Details    string         `json:"details,omitempty"`
	Internal   error          `json:"-"`
	StackTrace string         `json:"-"`
	Context    map[string]any `json:"context,omitempty"`
}

func (e *AppError) Error() string {
	if e.Internal != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Internal)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func (e *AppError) HTTPStatus() int {
	switch e.Code {
	case ErrCodeBadRequest, ErrCodeValidation, ErrCodeInvalidDocumentFormat, ErrCodeDocumentSizeExceeded:
		return http.StatusBadRequest
	case ErrCodeUnauthorized, ErrCodeInvalidCredentials, ErrCodeTokenExpired, ErrCodeTokenInvalid:
		return http.StatusUnauthorized
	case ErrCodeForbidden, ErrCodeDocumentAccessDenied, ErrCodeWorkspaceAccessDenied:
		return http.StatusForbidden
	case ErrCodeNotFound, ErrCodeUserNotFound, ErrCodeDocumentNotFound, ErrCodeWorkspaceNotFound, ErrCodeSessionNotFound:
		return http.StatusNotFound
	case ErrCodeConflict, ErrCodeUserAlreadyExists, ErrCodeDocumentLocked, ErrCodeOperationConflict, ErrCodeConcurrentEdit:
		return http.StatusConflict
	case ErrCodeTimeout:
		return http.StatusRequestTimeout
	case ErrCodeRateLimit:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}

func (e *AppError) WithContext(key string, value any) *AppError {
	if e.Context == nil {
		e.Context = make(map[string]any)
	}
	e.Context[key] = value
	return e
}

func (e *AppError) WithDetails(details string) *AppError {
	e.Details = details
	return e
}

func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StackTrace: getStackTrace(),
	}
}

func Wrap(err error, code ErrorCode, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		Internal:   err,
		StackTrace: getStackTrace(),
	}
}

func getStackTrace() string {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			buf = buf[:n]
			break
		}
		buf = make([]byte, 2*len(buf))
	}

	lines := strings.Split(string(buf), "\n")
	var filteredLines []string
	for i, line := range lines {
		if strings.Contains(line, "smart-document") && !strings.Contains(line, "errors.go") {
			filteredLines = append(filteredLines, line)
			if i+1 < len(lines) {
				filteredLines = append(filteredLines, lines[i+1])
			}
		}
	}
	return strings.Join(filteredLines, "\n")
}

// Common error constructors
func BadRequest(message string) *AppError {
	return New(ErrCodeBadRequest, message)
}

func Unauthorized(message string) *AppError {
	return New(ErrCodeUnauthorized, message)
}

func Forbidden(message string) *AppError {
	return New(ErrCodeForbidden, message)
}

func NotFound(message string) *AppError {
	return New(ErrCodeNotFound, message)
}

func Conflict(message string) *AppError {
	return New(ErrCodeConflict, message)
}

func Internal(err error, message string) *AppError {
	return Wrap(err, ErrCodeInternal, message)
}

func Validation(message string) *AppError {
	return New(ErrCodeValidation, message)
}

// Authentication errors
func InvalidCredentials() *AppError {
	return New(ErrCodeInvalidCredentials, "Invalid credentials provided")
}

func TokenExpired() *AppError {
	return New(ErrCodeTokenExpired, "Token has expired")
}

func TokenInvalid() *AppError {
	return New(ErrCodeTokenInvalid, "Invalid token provided")
}

func UserNotFound() *AppError {
	return New(ErrCodeUserNotFound, "User not found")
}

func UserAlreadyExists() *AppError {
	return New(ErrCodeUserAlreadyExists, "User already exists")
}

// Document errors
func DocumentNotFound() *AppError {
	return New(ErrCodeDocumentNotFound, "Document not found")
}

func DocumentAccessDenied() *AppError {
	return New(ErrCodeDocumentAccessDenied, "Access denied to document")
}

func DocumentLocked() *AppError {
	return New(ErrCodeDocumentLocked, "Document is locked for editing")
}

// Workspace errors
func WorkspaceNotFound() *AppError {
	return New(ErrCodeWorkspaceNotFound, "Workspace not found")
}

func WorkspaceAccessDenied() *AppError {
	return New(ErrCodeWorkspaceAccessDenied, "Access denied to workspace")
}

// Helper functions to check error types
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

func GetAppError(err error) (*AppError, bool) {
	appErr, ok := err.(*AppError)
	return appErr, ok
}

func HasCode(err error, code ErrorCode) bool {
	if appErr, ok := GetAppError(err); ok {
		return appErr.Code == code
	}
	return false
}

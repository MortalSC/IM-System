package errors

import "net/http"

// 示例错误
var (
	ErrForExample = NewErrEx("forExample", ErrCodeErrForExample, http.StatusForbidden, "forExample")
)

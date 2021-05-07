package api

import "fmt"

// Error ensures that protocol buffer errors implement error
func (e *Error) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

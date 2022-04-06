package routes

import (
	"encoding/json"
	"io"
	"net/http"
)

const (
	StatusSuccess = "success"
	StatusFail    = "fail"
	StatusError   = "error"
)

// Response data format for HTTP
type Response struct {
	Status  string      `json:"status" bson:"status"`   // Status code (error|fail|success)
	Code    int         `json:"code"  bson:"code"`      // HTTP status code
	Message string      `json:"message" bson:"message"` // Error or status message
	Data    interface{} `json:"data" bson:"data"`       // Data payload
}

// jsonHTTPEncode is a wrapper for json.NewEncoder(w).Encode(v)
func jsonHTTPEncode(w io.Writer, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

// convert the HTTP status code into JSend status
func getStatus(code int) string {
	if code >= 500 {
		return StatusError
	}
	if code >= 400 {
		return StatusFail
	}
	return StatusSuccess
}

// setHeaders set the default headers
func setHeaders(hw http.ResponseWriter, contentType string, code int) {
	hw.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	hw.Header().Set("Pragma", "no-cache")
	hw.Header().Set("Expires", "0")
	hw.Header().Set("Content-Type", contentType)
	hw.WriteHeader(code)
}

// sendResponse sends the HTTP response in JSON format
func sendResponse(hw http.ResponseWriter, hr *http.Request, code int, message string, data interface{}) {

	response := Response{
		Status:  getStatus(code),
		Code:    code,
		Message: message,
		Data:    data,
	}

	// send JSON response
	setHeaders(hw, "application/json", code)
	jsonHTTPEncode(hw, response)
}

package common

// HTTPResponse common data response for handling REST data and status
type HTTPResponse struct {
	ResponseCode    string `json:"responseCode,omitempty"`
	ResponseMessage string `json:"responseMessage,omitempty"`
	Data            any    `json:"data,omitempty"`
	HTTPStatus      int    `json:"-"`
}

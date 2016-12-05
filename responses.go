package stringsvc

// UppercaseResponse struct for Uppercase response
type UppercaseResponse struct {
    V   string `json:"v"`
    Err error  `json:"err"`
}

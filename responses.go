package stringsvc

// UppercaseResponse struct for Uppercase method
type UppercaseResponse struct {
    V   string `json:"v"`
    Err error  `json:"err"`
}

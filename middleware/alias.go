package middleware

import "github.com/l-vitaly/stringsvc"

// Alias Middleware for func(stringsvc.StringSvc) stringsvc.StringSvc
type Middleware func(stringsvc.StringSvc) stringsvc.StringSvc

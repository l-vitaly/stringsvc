package middleware

import "github.com/l-vitaly/stringsvc"

// Alias Middleware
type Middleware func(stringsvc.StringSvc) stringsvc.StringSvc

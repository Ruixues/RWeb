package RWeb

const (
	tMethodGet     = "GET"     // RFC 7231, 4.3.1
	tMethodHead    = "HEAD"    // RFC 7231, 4.3.2
	tMethodPost    = "POST"    // RFC 7231, 4.3.3
	tMethodPut     = "PUT"     // RFC 7231, 4.3.4
	tMethodPatch   = "PATCH"   // RFC 5789
	tMethodDelete  = "DELETE"  // RFC 7231, 4.3.5
	tMethodConnect = "CONNECT" // RFC 7231, 4.3.6
	tMethodOptions = "OPTIONS" // RFC 7231, 4.3.7
	tMethodTrace   = "TRACE"   // RFC 7231, 4.3.8
)

func ToRWebMethod(Method []byte) int {
	switch string(Method) {
	case tMethodGet:
		return MethodGet
	case tMethodHead:
		return MethodHead
	case tMethodPost:
		return MethodPost
	case tMethodPut:
		return MethodPut
	case tMethodPatch:
		return MethodPatch
	case tMethodDelete:
		return MethodDelete
	case tMethodConnect:
		return MethodConnect
	case tMethodOptions:
		return MethodOptions
	case tMethodTrace:
		return MethodTrace
	default:
		return MethodUnknown
	}
}

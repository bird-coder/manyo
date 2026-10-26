package errs

// StatusCode status code
//
//go:generate go run golang.org/x/tools/cmd/stringer -type=StatusCode -trimprefix Status
type StatusCode int

func (i StatusCode) Code() int {
	return int(i)
}

func (i StatusCode) Recoverable() bool {
	switch i {
	case StatusOptimisticLockError, StatusTryLockError:
		return true
	}
	return false
}

// Status like http code
const (
	StatusOK                           StatusCode = 0
	StatusUnknown                      StatusCode = 1
	StatusBadRequest                   StatusCode = 400 // RFC 9110, 15.5.1
	StatusUnauthorized                 StatusCode = 401 // RFC 9110, 15.5.2
	StatusPaymentRequired              StatusCode = 402 // RFC 9110, 15.5.3
	StatusForbidden                    StatusCode = 403 // RFC 9110, 15.5.4
	StatusNotFound                     StatusCode = 404 // RFC 9110, 15.5.5
	StatusMethodNotAllowed             StatusCode = 405 // RFC 9110, 15.5.6
	StatusNotAcceptable                StatusCode = 406 // RFC 9110, 15.5.7
	StatusProxyAuthRequired            StatusCode = 407 // RFC 9110, 15.5.8
	StatusRequestTimeout               StatusCode = 408 // RFC 9110, 15.5.9
	StatusConflict                     StatusCode = 409 // RFC 9110, 15.5.10
	StatusGone                         StatusCode = 410 // RFC 9110, 15.5.11
	StatusLengthRequired               StatusCode = 411 // RFC 9110, 15.5.12
	StatusPreconditionFailed           StatusCode = 412 // RFC 9110, 15.5.13
	StatusRequestEntityTooLarge        StatusCode = 413 // RFC 9110, 15.5.14
	StatusRequestURITooLong            StatusCode = 414 // RFC 9110, 15.5.15
	StatusUnsupportedMediaType         StatusCode = 415 // RFC 9110, 15.5.16
	StatusRequestedRangeNotSatisfiable StatusCode = 416 // RFC 9110, 15.5.17
	StatusExpectationFailed            StatusCode = 417 // RFC 9110, 15.5.18
	StatusTeapot                       StatusCode = 418 // RFC 9110, 15.5.19 (Unused)
	StatusMisdirectedRequest           StatusCode = 421 // RFC 9110, 15.5.20
	StatusUnprocessableEntity          StatusCode = 422 // RFC 9110, 15.5.21
	StatusLocked                       StatusCode = 423 // RFC 4918, 11.3
	StatusFailedDependency             StatusCode = 424 // RFC 4918, 11.4
	StatusTooEarly                     StatusCode = 425 // RFC 8470, 5.2.
	StatusUpgradeRequired              StatusCode = 426 // RFC 9110, 15.5.22
	StatusPreconditionRequired         StatusCode = 428 // RFC 6585, 3
	StatusTooManyRequests              StatusCode = 429 // RFC 6585, 4
	StatusRequestHeaderFieldsTooLarge  StatusCode = 431 // RFC 6585, 5
	StatusUnavailableForLegalReasons   StatusCode = 451 // RFC 7725, 3
	StatusInternalServerError          StatusCode = 500
)

// Internal
const (
	StatusServerMaintain    StatusCode = 600
	StatusNotEnoughResource StatusCode = 603
	StatusInstanceNotFound  StatusCode = 604
	StatusDuplicated        StatusCode = 605
)

// Recoverable
const (
	StatusOptimisticLockError StatusCode = 800
	StatusCacheDataError      StatusCode = 801 // 缓存数据错误
	StatusTryLockError        StatusCode = 802 // 尝试加锁失败
)

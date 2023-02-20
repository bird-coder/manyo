package redis

import "runtime/trace"

type traceConn struct {
	tr trace.Region
}

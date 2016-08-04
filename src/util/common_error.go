package util

import "errors"

var (
	ErrShutdown     = errors.New("connection is shut down")
	ErrRpcUnhandled = errors.New("rpc is unhandled")
)

const (
	RpcOk     = "ok"
	RpcFailed = "false"
)

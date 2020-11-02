package gservice

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"io"
)

type RpcType string

const (
	RpcTypeGRPC RpcType = "grpc"
)

func GRpcFactory() ConnFactory {
	return func(host string) (conn io.Closer, err error) {
		conn, err = grpc.Dial(host, grpc.WithInsecure())
		return
	}
}
func GRpcConnTest() ConnTestFunc {
	return func(closer io.Closer) error {
		conn := closer.(*grpc.ClientConn)
		status := conn.GetState()
		if status == connectivity.TransientFailure ||
			status == connectivity.Shutdown {
			return fmt.Errorf("conn status %s\n%w", status, ErrConn)
		} else {
			return nil
		}
	}
}

package bufconn

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufsize = 1024 * 1024

type Listener struct {
	sock *bufconn.Listener
}

func New() *Listener {
	return &Listener{
		sock: bufconn.Listen(bufsize),
	}
}

func (l *Listener) Sock() net.Listener {
	return l.sock
}

func (l *Listener) Close() error {
	return l.sock.Close()
}

func (l *Listener) Connect(ctx context.Context, opts ...grpc.DialOption) (cc *grpc.ClientConn, err error) {
	opts = append([]grpc.DialOption{grpc.WithContextDialer(l.Dialer)}, opts...)
	if cc, err = grpc.NewClient("passthrough://bufnet", opts...); err != nil {
		return nil, err
	}
	return cc, nil
}

func (l *Listener) Dialer(context.Context, string) (net.Conn, error) {
	return l.sock.Dial()
}

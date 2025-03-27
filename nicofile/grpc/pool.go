package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"sync"
	"time"
)

type Pool struct {
	conns   chan *grpc.ClientConn
	factory func() (*grpc.ClientConn, error)
	mu      sync.Mutex
	addr    string
}

func NewPool(addr string, poolSize int) (*Pool, error) {
	p := &Pool{
		conns: make(chan *grpc.ClientConn, poolSize),
		factory: func() (*grpc.ClientConn, error) {
			return grpc.NewClient(addr,
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithKeepaliveParams(keepalive.ClientParameters{
					Time:    30 * time.Second,
					Timeout: 10 * time.Second,
				}),
				grpc.WithInitialWindowSize(1<<24), // 16MB
				grpc.WithInitialConnWindowSize(1<<24),
				grpc.WithDefaultCallOptions(
					grpc.MaxCallRecvMsgSize(10<<20), // 10MB
					grpc.MaxCallSendMsgSize(10<<20),
				),
			)
		},
		addr: addr,
	}

	// 预热连接池
	for i := 0; i < poolSize; i++ {
		conn, err := p.factory()
		if err != nil {
			return nil, err
		}
		p.conns <- conn
	}

	return p, nil
}

func (p *Pool) Get() (*grpc.ClientConn, error) {
	select {
	case conn := <-p.conns:
		if p.check(conn) {
			return conn, nil
		}
		conn.Close()
		return p.factory()
	default:
		return p.factory()
	}
}

func (p *Pool) Put(conn *grpc.ClientConn) {
	if p.check(conn) {
		select {
		case p.conns <- conn:
		default:
			conn.Close()
		}
	} else {
		conn.Close()
	}
}

func (p *Pool) check(conn *grpc.ClientConn) bool {
	state := conn.GetState()
	return state != connectivity.Shutdown && state != connectivity.TransientFailure
}

func (p *Pool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	close(p.conns)
	for conn := range p.conns {
		conn.Close()
	}
}

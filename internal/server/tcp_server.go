package server

import (
	"bufio"
	"fmt"
	"goredis/common/config"
	"goredis/common/logger"
	"goredis/internal/protocol"
	"net"
	"time"
)

type (
	tcpserver struct {
		logger logger.Logger
		cfg    *config.Config
		ln     net.Listener
		exit   chan struct{}
		parser protocol.Parser
	}
)

var _ Server = (*tcpserver)(nil)

func NewTcpServer(opts ...ServerOption) *tcpserver {

	srvOptions := newServerOptions(opts...)
	return &tcpserver{
		cfg:    srvOptions.config,
		logger: srvOptions.logger,
		exit:   make(chan struct{}),
		parser: protocol.NewGrespParser(srvOptions.logger),
	}
}

func newServerOptions(opts ...ServerOption) *serverOptions {
	srvOpt := &serverOptions{}

	for _, opt := range opts {
		opt.apply(srvOpt)
	}

	return srvOpt
}

func (tsr *tcpserver) Start() error {

	ln, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", tsr.cfg.ServerOptions.Port))

	if err != nil {
		tsr.logger.Error(err)
		return err
	}

	tsr.ln = ln

	tsr.acceptLoop()
	return nil
}

func (tsr *tcpserver) acceptLoop() {
	for {
		conn, err := tsr.ln.Accept()
		if err != nil {
			select {
			case <-tsr.exit:
				return
			default:
				tsr.logger.Info(fmt.Sprintf("failed to established connection %s", err.Error()))
				continue
			}
		}

		tsr.logger.Info(fmt.Sprintf("Recieved connection from %s", conn.RemoteAddr().String()))
		go tsr.handleConn(conn)
	}
}

func (tsr *tcpserver) handleConn(conn net.Conn) {

	reader := bufio.NewReader(conn)

	for {

		err := tsr.parser.Parse(reader)
		if err != nil {
			tsr.logger.Error(err)
			break
		}

	}

}

func (tsr *tcpserver) Stop() {
	tsr.logger.Info("closing listner")
	tsr.ln.Close()
	close(tsr.exit)
	time.Sleep(10 * time.Second)
	tsr.logger.Info("go redis server completed shutdown")

}

package server

import (
	"bufio"
	"fmt"
	"goredis/common/config"
	"goredis/common/logger"
	"goredis/internal/event_processor"
	"goredis/internal/protocol"
	"io"
	"net"
	"time"
)

type (
	tcpserver struct {
		logger    logger.Logger
		cfg       *config.Config
		ln        net.Listener
		exit      chan struct{}
		parser    protocol.Parser
		eventLoop *event_processor.EventLoop
	}
)

var _ Server = (*tcpserver)(nil)

func NewTcpServer(opts ...ServerOption) *tcpserver {

	srvOptions := newServerOptions(opts...)
	return &tcpserver{
		cfg:       srvOptions.config,
		logger:    srvOptions.logger,
		exit:      make(chan struct{}),
		parser:    protocol.NewGrespParser(srvOptions.logger),
		eventLoop: event_processor.NewEventLoop(srvOptions.logger),
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

	tsr.logger.Info("starting event loop")
	tsr.eventLoop.Start()

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

		req, err := tsr.parser.Parse(reader)
		event := event_processor.NewEvent(conn)

		if err != nil && err != io.EOF {
			tsr.eventLoop.AddEvent(event.WithError(err))
			continue
		}

		if err == io.EOF {
			tsr.logger.Info(fmt.Sprintf("connection closed by client %s", conn.RemoteAddr()))
			conn.Close()
			break
		}

		tsr.eventLoop.AddEvent(event.WithRequest(req))
	}

}

func (tsr *tcpserver) Stop() {
	tsr.logger.Info("closing listner")
	tsr.ln.Close()
	close(tsr.exit)
	tsr.eventLoop.CloseLoop()
	time.Sleep(10 * time.Second)
	tsr.logger.Info("go redis server completed shutdown")

}

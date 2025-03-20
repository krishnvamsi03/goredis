package client

import (
	"bufio"
	"fmt"
	"goredis/internal/tokens"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

type (
	Client struct {
		address string
		port    int
		conn    net.Conn
	}

	ClientOpt func(*Client)
)

func WithAddress(address string) ClientOpt {
	return func(c *Client) {
		c.address = address
	}
}

func WithPort(port int) ClientOpt {
	return func(c *Client) {
		c.port = port
	}
}

func NewClient(clientOpts ...ClientOpt) *Client {

	client := &Client{}

	for _, opt := range clientOpts {
		opt(client)
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", client.address, client.port))
	if err != nil {
		panic(err)
	}
	client.conn = conn
	return client
}

func (cli *Client) Start() {

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	reader := bufio.NewReader(os.Stdin)
	conReader := bufio.NewReader(cli.conn)

	defer cli.conn.Close()

	go func() {
		<-quit
		cli.conn.Close()
		fmt.Print("quitting goredis bye!")
		os.Exit(0)
	}()

	fmt.Println("welcome to goredis, type quit to exit")
	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("error reading input, please try again")
			continue
		}

		if len(strings.TrimSpace(input)) == 0 {
			continue
		}

		if strings.TrimSpace(input) == "quit" {
			fmt.Println("quitting goredis bye!")
			break
		}

		cmd, err := cli.validateInp(input)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		requestBody := cli.buildRequest(cmd)
		_, err = cli.conn.Write([]byte(requestBody))
		if err != nil {
			fmt.Println(err)
			continue
		}
		response := ""
		contentLength := 0
		for i := 0; i < 2; i++ {
			res, err := conReader.ReadString('\n')
			if err != nil {
				fmt.Println(err)
				break
			}
			response += res
			cnt := strings.Split(res, " ")
			if len(cnt) >= 2 && strings.TrimSpace(cnt[0]) == tokens.CONTENT_LENGTH {
				contentLength, _ = strconv.Atoi(strings.TrimSpace(cnt[1]))
			}
		}

		content := ""
		for contentLength > 0 {
			ch, err := conReader.ReadByte()
			if err != nil {
				fmt.Println(err.Error())
				break
			}

			content += string(ch)
			contentLength -= 1
		}
		_, err = conReader.ReadByte()
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		response += content

		fmt.Println(cli.parseResponse(cmd.op, response))
	}
}

func (cli *Client) validateInp(inp string) (*command, error) {

	inps := strings.Split(inp, " ")
	for idx, value := range inps {
		inps[idx] = strings.TrimSpace(value)
	}

	switch inps[0] {
	case "PING":
		return &command{
			op: inps[0],
		}, nil
	case "GET":
		if len(inps) <= 1 {
			return nil, ErrInvalidArgs
		}

		return &command{
			op:  inps[0],
			key: inps[1],
		}, nil
	case "SET":
		if len(inps) <= 4 {
			return nil, ErrInvalidArgs
		}

		return &command{
			op:       inps[0],
			key:      inps[1],
			value:    inps[2],
			datatype: inps[3],
			ttl:      inps[4],
		}, nil
	case "DEL":
		if len(inps) <= 1 {
			return nil, ErrInvalidArgs
		}
		return &command{
			op:  inps[0],
			key: inps[1],
		}, nil
	case "PUSH":
		if len(inps) <= 2 {
			return nil, ErrInvalidArgs
		}

		return &command{
			op:    inps[0],
			key:   inps[1],
			value: inps[2],
		}, nil
	case "KEYS":
		if len(inps) <= 1 {
			return nil, ErrInvalidArgs
		}
		return &command{
			op:  inps[0],
			key: inps[1],
		}, nil
	case "POP":
		if len(inps) <= 1 {
			return nil, ErrInvalidCommand
		}
		cmd := &command{
			op:  inps[0],
			key: inps[1],
		}
		if len(inps) >= 3 {
			cmd.value = fmt.Sprintf("%s %s", inps[2], inps[3])
		}
		return cmd, nil
	case "INCR":
		if len(inps) <= 1 {
			return nil, ErrInvalidArgs
		}

		return &command{
			op:  inps[0],
			key: inps[1],
		}, nil
	case "DECR":
		if len(inps) <= 1 {
			return nil, ErrInvalidArgs
		}
		return &command{
			op:  inps[0],
			key: inps[1],
		}, nil
	default:
		return nil, ErrInvalidCommand
	}
}

func (cli *Client) buildRequest(cmd *command) string {

	switch cmd.op {
	case "PING":
		return "GRESP OP PING\n\n"
	case "GET":
		return fmt.Sprintf("GRESP OP GET KEY %s\n\n", cmd.key)
	case "SET":
		body := fmt.Sprintf("GRESP OP SET KEY %s DATA_TYPE %s TTL %s\n", cmd.key, cmd.datatype, cmd.ttl)
		body += fmt.Sprintf("CONTENT_LENGTH %d\n", len(cmd.value))
		body += fmt.Sprintf("%s\n\n", cmd.value)
		return body
	case "DEL":
		return fmt.Sprintf("GRESP OP DEL KEY %s\n\n", cmd.key)
	case "PUSH":
		body := fmt.Sprintf("GRESP OP PUSH KEY %s\n", cmd.key)
		body += fmt.Sprintf("CONTENT_LENGTH %d\n", len(cmd.value))
		body += fmt.Sprintf("%s\n\n", cmd.value)
		return body
	case "POP":
		body := fmt.Sprintf("GRESP OP POP KEY %s\n", cmd.key)
		if len(cmd.value) > 0 {
			body += fmt.Sprintf("CONTENT_LENGTH %d\n", len(cmd.value))
			body += fmt.Sprintf("%s\n", cmd.value)
		}
		body += "\n"
		return body
	case "INCR":
		return fmt.Sprintf("GRESP OP INCR KEY %s\n\n", cmd.key)
	case "DECR":
		return fmt.Sprintf("GRESP OP DECR KEY %s \n\n", cmd.key)
	case "KEYS":
		return fmt.Sprintf("GRESP OP KEYS KEY %s\n\n", cmd.key)
	default:
		return ""
	}

}

func (cli *Client) parseResponse(op, response string) string {

	lines := strings.Split(response, "\n")
	if len(lines) < 3 {
		return "invalid response received from server"
	}

	meta := strings.Split(lines[0], " ")
	dt := meta[2]

	switch dt {
	case "LIST":
		res := strings.Split(lines[2], ":")
		return "[" + strings.Join(res, ",") + "]"
	default:
		switch op {
		case "KEYS":
			ans := ""
			for _, value := range lines[2:] {
				ans += fmt.Sprintf("%s\n", value)
			}

			return strings.TrimSpace(ans)
		default:
			return lines[2]
		}
	}
}

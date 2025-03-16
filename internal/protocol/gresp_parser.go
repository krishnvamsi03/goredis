package protocol

import (
	"bufio"
	"errors"
	"goredis/common/logger"
	gerrors "goredis/errors"
	"goredis/internal/request"
	"goredis/internal/tokens"
	"strconv"
	"strings"
)

type (
	grespProtocolParser struct {
		logger logger.Logger
	}
)

var _ Parser = (*grespProtocolParser)(nil)

func NewGrespParser(logger logger.Logger) *grespProtocolParser {
	return &grespProtocolParser{
		logger: logger,
	}
}

func (grp *grespProtocolParser) Parse(reader *bufio.Reader) (*request.Request, error) {

	commands := []string{}

	line, err := reader.ReadString('\n')
	line = strings.TrimSpace(line)

	if err != nil {
		return nil, err
	}

	// first line should not be empty
	if len(line) == 0 {
		return nil, gerrors.ErrIncomReqBody
	}

	commands = append(commands, strings.Split(line, " ")...)

	err = grp.validateCommands(commands)
	if err != nil {
		return nil, err
	}

	cmd, err := grp.getCommand(commands)
	if err != nil {
		grp.logger.Error(err)
		return nil, err
	}

	err = grp.readContentIfExists(cmd, reader)
	if err != nil {
		grp.logger.Error(err)
		return nil, err
	}

	return cmd, nil
}

func (grp *grespProtocolParser) validateCommands(commands []string) error {

	cnt := 0
	for _, cmd := range commands {
		if _, ok := tokens.MANDATORY_TOKENS[cmd]; ok {
			cnt += 1
		}
	}

	if cnt < len(tokens.MANDATORY_TOKENS) {
		return gerrors.ErrIncomReqBody
	}

	return nil
}

func (grp *grespProtocolParser) getCommand(commands []string) (*request.Request, error) {

	commands = commands[1:]

	i := 0
	cmds := []request.RequestOptions{}

	for i < len(commands) {
		cmd := commands[i]
		cmd = strings.TrimSpace(cmd)
		if len(cmd) == 0 {
			i++
			continue
		}

		i += 1
		if f, ok := tokens.COMMANDS[cmd]; ok && i < len(commands) {
			cmds = append(cmds, f(commands[i]))
			i += 1
		}
	}

	cmd := &request.Request{}

	for _, cd := range cmds {
		cd.Apply(cmd)
	}

	err := cmd.Validate()
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

func (grp *grespProtocolParser) readContentIfExists(req *request.Request, reader *bufio.Reader) error {

	switch strings.ToLower(*req.Op) {
	case tokens.SET.ToLower(), tokens.PUSH.ToLower():
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		line = strings.TrimSpace(line)
		value, err := grp.readContent(line, reader)
		if err != nil {
			return err
		}
		req.Value = value
	case tokens.POP.ToLower():
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		line = strings.TrimSpace(line)
		emptyValue := ""
		if len(line) == 0 {
			req.Value = &emptyValue
			return nil
		}

		value, err := grp.readContent(line, reader)
		if err != nil {
			return err
		}
		req.Value = value

	default:
		emptyLine, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		if len(strings.TrimSpace(emptyLine)) > 0 {
			return gerrors.ErrInvalidProtocol
		}
	}
	return nil
}

func (grp *grespProtocolParser) readContent(line string, reader *bufio.Reader) (*string, error) {

	words := strings.Split(line, " ")
	if len(words) <= 1 {
		return nil, gerrors.ErrInvalidContentType
	}
	if words[0] != tokens.CONTENT_LENGTH {
		return nil, gerrors.ErrInvalidProtocol
	}

	contentLen, err := strconv.Atoi(words[1])
	if err != nil {
		return nil, errors.Join(gerrors.ErrInvalidProtocol, err)
	}

	if contentLen == 0 {
		return nil, gerrors.ErrInvalidContentType
	}

	i := 0
	value := ""
	for i < contentLen {
		c, err := reader.ReadByte()
		if err != nil {
			return nil, gerrors.ErrContentMismatch
		}

		value += string(c)
		i += 1
	}
	_, err = reader.ReadByte()
	if err != nil {
		return nil, gerrors.ErrInvalidContentType
	}
	return &value, nil
}

package protocol

import (
	"bufio"
	"fmt"
	"goredis/common/logger"
	gerrors "goredis/errors"
	"goredis/internal/tokens"
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

func (grp *grespProtocolParser) IsStart(line string) error {
	tkns := strings.Split(line, " ")
	cnt := 0
	for _, token := range tkns {
		if _, ok := tokens.INIT_TOKENS[token]; ok {
			cnt += 1
		}
	}

	if cnt < len(tokens.INIT_TOKENS) {
		return gerrors.ErrInvalidProtocol
	}

	return nil
}

func (grp *grespProtocolParser) Parse(reader *bufio.Reader) error {

	actions := []string{}

	for i := 0; i < 2; i++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		// first two lines should not be empty
		if len(strings.TrimSpace(line)) == 0 {
			return gerrors.ErrIncomReqBody
		}

		actions = append(actions, strings.Split(line, " ")...)
	}

	fmt.Println(actions)
	return nil
}

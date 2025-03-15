package command

import gerrors "goredis/errors"

type (
	Command struct {
		Key      *string
		Op       *string
		Expr     *string
		Value    *string
		Datatype *string
	}

	CommandOptions interface {
		Apply(*Command)
	}

	CommandApplyFunc func(*Command)
)

func (c *Command) Validate() error {
	if c.Key == nil || c.Op == nil || c.Datatype == nil {
		return gerrors.ErrInvalidCmds
	}
	return nil
}

func (f CommandApplyFunc) Apply(command *Command) {
	f(command)
}

func WithOp(operation string) CommandOptions {
	return CommandApplyFunc(func(c *Command) {
		c.Op = &operation
	})
}

func WithKey(key string) CommandOptions {
	return CommandApplyFunc(func(c *Command) {
		c.Key = &key
	})
}

func WithExpr(expr string) CommandOptions {
	return CommandApplyFunc(func(c *Command) {
		c.Expr = &expr
	})
}

func WithValue(value string) CommandOptions {
	return CommandApplyFunc(func(c *Command) {
		c.Value = &value
	})
}

func WithDataType(dt string) CommandOptions {
	return CommandApplyFunc(func(c *Command) {
		c.Datatype = &dt
	})
}

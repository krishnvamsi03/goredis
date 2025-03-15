package request

import gerrors "goredis/errors"

type (
	Request struct {
		Key      *string
		Op       *string
		Expr     *string
		Value    *string
		Datatype *string
	}

	RequestOptions interface {
		Apply(*Request)
	}

	RequestApplyFunc func(*Request)
)

func (req *Request) Validate() error {
	if req.Key == nil || req.Op == nil || req.Datatype == nil {
		return gerrors.ErrInvalidCmds
	}
	return nil
}

func (rf RequestApplyFunc) Apply(r *Request) { rf(r) }

func WithOp(operation string) RequestOptions {
	return RequestApplyFunc(func(c *Request) {
		c.Op = &operation
	})
}

func WithKey(key string) RequestOptions {
	return RequestApplyFunc(func(c *Request) {
		c.Key = &key
	})
}

func WithExpr(expr string) RequestOptions {
	return RequestApplyFunc(func(c *Request) {
		c.Expr = &expr
	})
}

func WithValue(value string) RequestOptions {
	return RequestApplyFunc(func(c *Request) {
		c.Value = &value
	})
}

func WithDataType(dt string) RequestOptions {
	return RequestApplyFunc(func(c *Request) {
		c.Datatype = &dt
	})
}

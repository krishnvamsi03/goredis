package response

import (
	"fmt"
	"goredis/internal/utils"
)

type (
	Response struct {
		code     string
		isOk     bool
		res      string
		datatype string
	}
)

func NewResponse() *Response {
	return &Response{}
}

func (res *Response) WithCode(code string) *Response {
	res.code = code
	return res
}

func (res *Response) WithOk(isOk bool) *Response {
	res.isOk = isOk
	return res
}

func (res *Response) WithRes(result string) *Response {
	res.res = result
	return res
}

func (res *Response) WithDatatype(datatype string) *Response {
	res.datatype = datatype
	return res
}

func (res *Response) Build() string {
	isOk := "NOT_OK"
	if res.isOk {
		isOk = "OK"
	}

	defaultDt := "STR"
	if !utils.IsEmpty(res.datatype) {
		defaultDt = res.datatype
	}

	response := fmt.Sprintf("GRESP %s %s %s\nCONTENT_LENGTH %d\n%s\n", isOk, defaultDt, res.code, len(res.res), res.res)
	return response
}

func BuildResponseWithError(err error) string {
	return "GRESP NOT_OK\n" + fmt.Sprintf("CONTENT_LENGTH %d\n", len(err.Error())) + fmt.Sprintf("%s\n", err.Error())
}

func BuildResponseWithMsg(value string) string {
	return "GRESP OK\n" + fmt.Sprintf("CONTENT_LENGTH %d\n", len(value)) + fmt.Sprintf("%s\n", value)
}

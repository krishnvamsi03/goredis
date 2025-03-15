package response

import "fmt"

func BuildResponseWithError(err error) string {
	return "GRESP NOT_OK\n" + fmt.Sprintf("CONTENT_LENGTH %d\n", len(err.Error())) + fmt.Sprintf("%s\n", err.Error())
}

func BuildResponseWithMsg(value string) string {
	return "GRESP OK\n" + fmt.Sprintf("CONTENT_LENGTH %d\n", len(value)) + fmt.Sprintf("%s\n", value)
}

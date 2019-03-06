package sunshinemotion

import "fmt"

// see IsServiceError()
type serviceError struct {
	status  int64
	message string
}

func (e serviceError) Error() string {
	return fmt.Sprintf(`Service Status %d, message "%s"`, e.status, e.message)
}

type networkError struct {
	message string
}

func (e networkError) Error() string {
	return fmt.Sprintf(`Network Error: %s`, e.message)
}

// serviceError指的是：解析出有效结果以后，在结果中，RPC服务器明确表示的错误，属于业务错误。
// 通常是预期之中的，如传入的帐号信息无效。
func IsServiceError(err error) bool {
	_, ok := err.(serviceError)
	return ok
}

// networkError指的是：从本地调用到RPC服务器接受到请求
// 或者RPC服务器对本地调用进行响应过程中发生的错误
func IsNetworkError(err error) bool {
	_, ok := err.(networkError)
	return ok
}

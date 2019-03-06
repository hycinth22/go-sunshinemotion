package sunshinemotion

var serviceErrorLookupMap = make(map[int64]error)

func init() {
	serviceErrorLookupMap[1] = nil
	serviceErrorLookupMap[2] = ErrTokenExpired
}

// 将code和message转化为go的标准错误
// 优先根据code转化为已存在的错误，否则动态生成一个错误
func serviceCodeToGoError(code int64, message string) error {
	err, ok := serviceErrorLookupMap[code]
	if ok {
		return err
	}
	return serviceError{code, message}
}

package ssmt

var (
	serviceErrorTable          = make(map[int64]IServiceError)
	ErrWrongUsernameOrPassword = ServiceError{0, "用户名或者密码错误", ""}
	ErrTokenExpired            = ServiceError{2, "超时，请重新登录", ""}
	ErrDisqualifiedSpeed       = ServiceError{5, "速度不合格", ""}
)

func init() {
	serviceErrorTable[0] = ErrWrongUsernameOrPassword
	serviceErrorTable[1] = nil
	serviceErrorTable[2] = ErrTokenExpired
	serviceErrorTable[5] = ErrDisqualifiedSpeed
}

func translateServiceError(statusCode int64, statusMessage string) IServiceError {
	err, exist := serviceErrorTable[statusCode]
	if exist {
		if err != nil {
			err.SetMsg(statusMessage)
		}
		return err
	}
	return ServiceError{statusCode, "Unknown", statusMessage}
}

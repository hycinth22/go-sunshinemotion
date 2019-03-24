package ssmt

var (
	serviceErrorTable          = make(map[int64]IServiceError)
	ErrForbidden               = ServiceError{-1, "无权进行此操作，可能不在跑步时间或被封禁", ""}
	ErrWrongUsernameOrPassword = ServiceError{0, "用户名或者密码错误", ""}
	ErrTokenExpired            = ServiceError{2, "超时，请重新登录", ""}
	ErrDisqualifiedSpeed       = ServiceError{5, "速度不合格", ""}
	ErrIllegalData             = ServiceError{100001, "非法数据", ""}
)

func init() {
	serviceErrorTable[-1] = ErrForbidden
	serviceErrorTable[0] = ErrWrongUsernameOrPassword
	serviceErrorTable[1] = nil
	serviceErrorTable[2] = ErrTokenExpired
	serviceErrorTable[5] = ErrDisqualifiedSpeed
	serviceErrorTable[100001] = ErrIllegalData
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

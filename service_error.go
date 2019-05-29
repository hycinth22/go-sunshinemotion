package ssmt

var (
	serviceErrorTable    = make(map[int64]IServiceError)
	ErrForbidden         = ServiceError{-1, "禁止进行此操作", ""} // 无权进行此操作，可能不在跑步时间或帐号被封禁
	ErrInvalidIdentity   = ServiceError{0, "身份信息错误", ""}
	ErrInvalidToken      = ServiceError{2, "令牌信息无效", ""}
	ErrDisqualifiedSpeed = ServiceError{5, "速度不合格", ""}
	ErrIllegalData       = ServiceError{100001, "非法数据", ""}

	ErrTokenExpired            = ServiceError{2, "无效的登录信息", ""}            // DEPRECATED
	ErrWrongUsernameOrPassword = ServiceError{0, "用户名或者密码错误，或体能测试未开放", ""} // DEPRECATED 用户名或者密码错误，或体能测试未开放
)

func init() {
	serviceErrorTable[-1] = ErrForbidden
	serviceErrorTable[0] = ErrWrongUsernameOrPassword
	serviceErrorTable[1] = nil
	serviceErrorTable[2] = ErrInvalidToken
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

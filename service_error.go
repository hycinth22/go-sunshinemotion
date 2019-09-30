package ssmt

var (
	serviceErrorTable    = make(map[int64]IServiceError)
	ErrForbidden         = &ServiceError{-1, "禁止进行此操作", ""} // 无权进行此操作，如：不在跑步时间、帐号被封禁
	ErrInvalidIdentity   = &ServiceError{0, "身份信息错误", ""}   // 当前身份信息不足以完成该操作。如：提供的用户名或者密码错误、体能测试未开放
	ErrInvalidToken      = &ServiceError{2, "令牌信息无效", ""}
	ErrDisqualifiedSpeed = &ServiceError{5, "速度不合格", ""}
	ErrIllegalData       = &ServiceError{100001, "非法数据", ""}

	ErrTokenExpired            = &ServiceError{2, "无效的登录信息", ""}            // DEPRECATED
	ErrWrongUsernameOrPassword = &ServiceError{0, "用户名或者密码错误，或体能测试未开放", ""} // DEPRECATED: 用户名或者密码错误，或体能测试未开放
)

func init() {
	serviceErrorTable[1] = nil
	registerServiceErrors(ErrForbidden, ErrWrongUsernameOrPassword, ErrInvalidToken, ErrDisqualifiedSpeed, ErrIllegalData)
}

func registerServiceErrors(errors ...IServiceError) {
	var ok bool
	for _, err := range errors {
		ok = registerServiceError(err, false)
		if !ok {
			panic("Register IServiceError: Duplicated Status Code")
		}
	}
}

func registerServiceError(err IServiceError, overwrite bool) bool {
	code := err.GetCode()
	if !overwrite {
		_, exist := serviceErrorTable[code]
		if exist {
			return false
		}
	}
	serviceErrorTable[code] = err
	return true
}

func translateServiceError(statusCode int64, statusMessage string) IServiceError {
	err, exist := serviceErrorTable[statusCode]
	if exist {
		if err != nil {
			err.SetMsg(statusMessage)
		}
		return err
	}
	return &ServiceError{statusCode, "Unknown", statusMessage}
}

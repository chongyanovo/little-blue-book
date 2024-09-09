package ginx

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func Success() Result {
	return Result{
		Code: 0,
		Msg:  "请求成功",
		Data: nil,
	}
}

func SuccessWithMsg(msg string) Result {
	return Result{
		Code: 0,
		Msg:  msg,
		Data: nil,
	}
}

func SuccessWithData(msg string, data any) Result {
	return Result{
		Code: 0,
		Msg:  msg,
		Data: data,
	}
}

func Fail() Result {
	return Result{
		Code: 1,
		Msg:  "请求失败",
		Data: nil,
	}
}
func FailWithMsg(msg string) Result {
	return Result{
		Code: 1,
		Msg:  msg,
		Data: nil,
	}
}

func FailWithData(msg string, data any) Result {
	return Result{
		Code: 1,
		Msg:  msg,
		Data: data,
	}
}

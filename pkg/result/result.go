package result

import (
	"encoding/json"
)

type Response struct {
	Code int         `json:"code"` // 错误码
	Msg  string      `json:"msg"`  // 错误描述
	Data interface{} `json:"data"` // 返回数据
}

// 自定义响应信息
func (res *Response) WithMsg(message string) Response {
	return Response{
		Code: res.Code,
		Msg:  message,
		Data: res.Data,
	}
}

// 追加响应数据
func (res *Response) WithData(data interface{}) Response {
	return Response{
		Code: res.Code,
		Msg:  res.Msg,
		Data: data,
	}
}

// ToString 返回 JSON 格式的错误详情
func (res *Response) ToString() string {
	err := &struct {
		Code int         `json:"code"`
		Msg  string      `json:"msg"`
		Data interface{} `json:"data"`
	}{
		Code: res.Code,
		Msg:  res.Msg,
		Data: res.Data,
	}
	raw, _ := json.Marshal(err)
	return string(raw)
}

// 构造函数
func response(code int, msg string) *Response {
	return &Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
}

// 错误码规则:
// (1) 错误码需为 > 0 的数;
//
// (2) 错误码为 5 位数:
//              ----------------------------------------------------------
//                  第1位               2、3位                  4、5位
//              ----------------------------------------------------------
//                服务级错误码          模块级错误码	         具体错误码
//              ----------------------------------------------------------

var (
	// OK
	OK  = response(200, "ok") // 通用成功
	Err = response(500, "")   // 通用错误

	// 服务级错误码
	ErrParam     = response(10001, "参数有误")
	ErrImageParam = response(10002, "镜像参数不能为空")
	ErrAppNameParam = response(10003, "应用名参数不能为空")
	ErrNamespaceParam = response(10003, "Namespace不能为空")

	// 模块级错误码 - deployment
	ErrDeployment = response(20100, "用户服务异常")
	ErrDeploymentCreate   = response(20101, "创建deployment错误")
	ErrDeploymentFinished = response(20102, "Deployment发布失败")
	ErrDeploymentAlreadyExist = response(20103, "Deployment已存在")
	ErrDeploymentNotFound = response(20104, "Deployment不存在")
)



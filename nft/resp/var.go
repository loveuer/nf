package resp

const (
	MSG200 = "请求成功"
	MSG202 = "请求成功, 请稍后..."
	MSG400 = "请求参数错误"
	MSG401 = "登录已过期, 请重新登录"
	MSG403 = "请求权限不足"
	MSG404 = "请求资源未找到"
	MSG418 = "请求条件不满足, 请稍后再试"
	MSG429 = "请求过于频繁, 请稍后再试"
	MSG500 = "服务器开小差了, 请稍后再试"
	MSG501 = "功能开发中, 尽情期待"
)

const (
	RealStatusHeader = "NF-STATUS"
)

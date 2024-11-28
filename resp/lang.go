package resp

// 定义 Lang 类型
type Lang string

// 定义语言常量
const (
	Zh Lang = "zh"
	En Lang = "en"
	Es Lang = "pt"
	Th Lang = "th"
)

// 定义语言列表变量，便于后续扩展
var Langs = []Lang{Zh, En, Es, Th}

// 获取语言字符串的方法
func (a Lang) String() string {
	return string(a)
}

package fetch

// Status 状态
type Status struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// Protocol 协议层
type Protocol struct {
	Status
	Data interface{} `json:"data"`
}

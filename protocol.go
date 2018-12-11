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

// Get send a Get Request
func (p *Protocol) Get(url string, option *RequestData) error {
	return request(GET, url, p, option)
}

func (p *Protocol) Post(url string, option *RequestData) error {
	return request(POST, url, p, option)
}

// NewProtocol 产生新的Protocol
func NewProtocol(data interface{}) *Protocol {
	return &Protocol{
		Data: data,
	}
}

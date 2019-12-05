package sakura

type Sakura struct {
}

func NewSakura() *Sakura {
	return &Sakura{}
}

func (s Sakura) Version() string {
	return "0.1.x"
}

func (s Sakura) Author() string {
	return "zhenhaowang"
}

func (s Sakura) Email() string {
	return "hit.zhenhao@gmail.com"
}

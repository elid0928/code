package factory

type Human interface {
	// Name string
	Skin() string
	Language() string
}

func NewAsia() *asia {
	return &asia{}
}

type asia struct {
}

func (a *asia) Skin() string {
	return "yellow"
}

func (a *asia) Language() string {
	return "Chinese"

}

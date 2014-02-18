package spaghetti

type Message struct {
	ConnId string `json:"-"`
	Body   interface{}
}

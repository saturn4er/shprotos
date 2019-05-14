package shprotos

type Service struct {
	Name          string
	QuotedComment string
	Methods       []*Method
	File          *File
}

type Method struct {
	Name          string
	QuotedComment string
	InputMessage  *Message
	OutputMessage *Message
	Service       *Service
}

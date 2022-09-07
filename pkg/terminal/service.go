package terminal

type Size struct {
	Cols   int `json:"cols"`
	Rows   int `json:"rows"`
	Width  int `json:"width"`
	Height int `json:"height"`
}
type Service interface {
	Start(size *Size) (err error)
	Write(buf []byte) (n int, err error)
	Read(buf []byte) (n int, err error)
	ChangeSize(size *Size) (err error)
	Stop()
}

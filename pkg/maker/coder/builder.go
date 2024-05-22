package coder

import "os"

func (this_ *Coder) NewBuilder(path string) (res *Builder, err error) {
	f, err := this_.CreateAndOpen(path)
	if err != nil {
		return
	}
	res = &Builder{
		path:      path,
		f:         f,
		rowNumber: 1,
		colNumber: 1,
	}
	return
}

type Builder struct {
	path      string
	f         *os.File
	tab       int
	rowNumber int
	colNumber int
}

func (this_ *Builder) AppendLine(line string) (err error) {
	str := ""

	for i := 0; i < this_.tab; i++ {
		str += "    "
	}
	str += line
	str += "\n"
	_, err = this_.f.WriteString(str)
	return
}

func (this_ *Builder) NewLine() (err error) {
	_, err = this_.f.WriteString("\n")
	return
}

func (this_ *Builder) Close() {
	_ = this_.f.Close()
}

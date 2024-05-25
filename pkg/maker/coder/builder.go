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

func (this_ *Builder) GetRowNumber() int {
	return this_.rowNumber
}
func (this_ *Builder) GetColNumber() int {
	return this_.colNumber
}
func (this_ *Builder) Tab() {
	this_.tab++
}
func (this_ *Builder) Indent() {
	if this_.tab > 0 {
		this_.tab--
	}
}
func (this_ *Builder) SetTab(tab int) {
	this_.tab = tab
}

func (this_ *Builder) GetTab() (tab int) {
	tab = this_.tab
	return
}

func (this_ *Builder) AppendTabLine(ss ...string) {
	str := ""

	for i := 0; i < this_.tab; i++ {
		str += "    "
	}
	for _, s := range ss {
		str += s
	}
	str += "\n"
	_, _ = this_.f.WriteString(str)
	this_.rowNumber++
	return
}

func (this_ *Builder) AppendTab() {
	str := ""
	for i := 0; i < this_.tab; i++ {
		str += "    "
	}
	_, _ = this_.f.WriteString(str)
	return
}

func (this_ *Builder) AppendCode(ss ...string) {
	str := ""
	for _, s := range ss {
		str += s
	}
	_, _ = this_.f.WriteString(str)
	return
}

func (this_ *Builder) NewLine() {
	_, _ = this_.f.WriteString("\n")
	this_.rowNumber++
	return
}

func (this_ *Builder) Close() {
	_ = this_.f.Close()
}

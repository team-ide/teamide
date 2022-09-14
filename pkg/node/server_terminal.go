package node

import (
	"teamide/pkg/terminal"
	"teamide/pkg/util"
)

func (this_ *Server) TerminalStart(lineNodeIdList []string, size *terminal.Size, onRead func(buf []byte) (err error), onReadError func(buf []byte) (err error)) (key string, err error) {
	readKey := util.UUID()
	readErrorKey := util.UUID()
	this_.addOnBytesCache(readKey, &OnBytes{
		start: func() (err error) {
			return
		},
		on: func(buf []byte) (err error) {
			err = onRead(buf)
			return
		},
		end: func() (err error) {
			return
		},
	})

	this_.addOnBytesCache(readErrorKey, &OnBytes{
		start: func() (err error) {
			return
		},
		on: func(buf []byte) (err error) {
			err = onReadError(buf)
			return
		},
		end: func() (err error) {
			return
		},
	})

	key, err = this_.workTerminalStart(lineNodeIdList, size, readKey, readErrorKey)
	if err != nil {
		return
	}
	return
}

func (this_ *Server) TerminalWrite(lineNodeIdList []string, key string, buf []byte) (err error) {

	err = this_.workTerminalWrite(lineNodeIdList, key, buf)
	if err != nil {
		return
	}
	return
}

func (this_ *Server) TerminalChangeSize(lineNodeIdList []string, key string, size *terminal.Size) (err error) {

	err = this_.workTerminalChangeSize(lineNodeIdList, key, size)
	if err != nil {
		return
	}
	return
}

func (this_ *Server) TerminalIsWindows(lineNodeIdList []string) (isWindows bool, err error) {
	isWindows, err = this_.workTerminalIsWindows(lineNodeIdList)
	if err != nil {
		return
	}
	return
}

func (this_ *Server) TerminalStop(lineNodeIdList []string, key string) (err error) {

	err = this_.workTerminalStop(lineNodeIdList, key)
	if err != nil {
		return
	}
	return
}

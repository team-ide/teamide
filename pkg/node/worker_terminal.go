package node

import (
	"errors"
	"go.uber.org/zap"
	"teamide/pkg/terminal"
	"teamide/pkg/util"
)

func (this_ *Worker) workTerminalStart(lineNodeIdList []string, size *terminal.Size, readKey string, readErrorKey string) (key string, err error) {
	send, err := this_.sendToNext(lineNodeIdList, "", func(listener *MessageListener) (e error) {
		res, e := this_.Call(listener, methodTerminalStart, &Message{
			LineNodeIdList: lineNodeIdList,
			TerminalWorkData: &TerminalWorkData{
				Size:         size,
				ReadKey:      readKey,
				ReadErrorKey: readErrorKey,
			},
		})
		if e != nil {
			return
		}

		if res != nil && res.TerminalWorkData != nil {
			key = res.TerminalWorkData.Key
		}

		return
	})
	if err != nil || send {
		return
	}

	service := terminal.NewLocalService()
	err = service.Start(size)
	if err != nil {
		return
	}

	Logger.Info("local service start success")

	key = util.UUID()

	this_.addTerminalService(key, service)

	var line []string
	for i := len(lineNodeIdList) - 1; i >= 0; i-- {
		line = append(line, lineNodeIdList[i])
	}
	go func() {
		defer func() {
			this_.removeTerminalService(key)
			service.Stop()
			Logger.Info("local service stopped")
		}()
		err = this_.workSend(line, readKey, service.Read)
		if err != nil {
			Logger.Error("terminal read send error", zap.Error(err))
		}
	}()

	go func() {
		defer func() {
			this_.removeTerminalService(key)
			service.Stop()
			Logger.Info("local service stopped")
		}()
		err = this_.workSend(line, readErrorKey, service.ReadError)
		if err != nil {
			Logger.Error("terminal read error send error", zap.Error(err))
		}
	}()

	return
}

func (this_ *Worker) workTerminalChangeSize(lineNodeIdList []string, key string, size *terminal.Size) (err error) {
	send, err := this_.sendToNext(lineNodeIdList, "", func(listener *MessageListener) (e error) {
		_, e = this_.Call(listener, methodTerminalChangeSize, &Message{
			LineNodeIdList: lineNodeIdList,
			TerminalWorkData: &TerminalWorkData{
				Key:  key,
				Size: size,
			},
		})
		if e != nil {
			return
		}

		return
	})
	if err != nil || send {
		return
	}

	service := this_.getTerminalService(key)
	if service == nil {
		err = errors.New("service [" + key + "] is not found.")
		return
	}

	err = service.ChangeSize(size)

	return
}

func (this_ *Worker) workTerminalStop(lineNodeIdList []string, key string) (err error) {
	send, err := this_.sendToNext(lineNodeIdList, "", func(listener *MessageListener) (e error) {
		_, e = this_.Call(listener, methodTerminalStop, &Message{
			LineNodeIdList: lineNodeIdList,
			TerminalWorkData: &TerminalWorkData{
				Key: key,
			},
		})
		if e != nil {
			return
		}

		return
	})
	if err != nil || send {
		return
	}

	service := this_.getTerminalService(key)
	if service == nil {
		err = errors.New("service [" + key + "] is not found.")
		return
	}

	service.Stop()

	return
}

func (this_ *Worker) workTerminalIsWindows(lineNodeIdList []string) (isWindows bool, err error) {
	send, err := this_.sendToNext(lineNodeIdList, "", func(listener *MessageListener) (e error) {
		res, e := this_.Call(listener, methodTerminalIsWindows, &Message{
			LineNodeIdList:   lineNodeIdList,
			TerminalWorkData: &TerminalWorkData{},
		})
		if e != nil {
			return
		}

		if res != nil && res.TerminalWorkData != nil {
			isWindows = res.TerminalWorkData.IsWindows
		}

		return
	})
	if err != nil || send {
		return
	}

	service := terminal.NewLocalService()
	isWindows, err = service.IsWindows()

	return
}

func (this_ *Worker) workTerminalWrite(lineNodeIdList []string, key string, buf []byte) (err error) {
	send, err := this_.sendToNext(lineNodeIdList, key, func(listener *MessageListener) (e error) {
		_, e = this_.Call(listener, methodTerminalWrite, &Message{
			LineNodeIdList: lineNodeIdList,
			Bytes:          buf,
			HasBytes:       true,
			TerminalWorkData: &TerminalWorkData{
				Key: key,
			},
		})
		if e != nil {
			return
		}

		return
	})
	if err != nil || send {
		return
	}

	service := this_.getTerminalService(key)
	if service == nil {
		err = errors.New("service [" + key + "] is not found.")
		return
	}

	_, err = service.Write(buf)

	return
}

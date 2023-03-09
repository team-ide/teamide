package module

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"teamide/internal/context"
	"teamide/pkg/base"
)

type ListenResponse struct {
	Events []*context.ListenEvent `json:"events"`
}

func (this_ *Api) listen(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
	}()
	var userId int64
	var clientTabKey = requestBean.ClientTabKey
	var clientKey = requestBean.ClientKey
	if requestBean.JWT != nil {
		userId = requestBean.JWT.UserId
	}
	if clientTabKey == "" {
		err = errors.New("client tab key is null")
		return
	}
	if clientKey == "" {
		err = errors.New("client tab key is null")
		return
	}
	listener := context.GetListener(clientTabKey)
	if listener == nil {
		listener = context.NewClientTabListener(clientKey, clientTabKey, userId)
		context.AddListener(listener)
	}
	if listener.ClientKey != clientKey {
		context.ChangeListenerClientKey(listener, clientKey)
	}
	if listener.UserId != userId {
		context.ChangeListenerUserId(listener, userId)
	}

	response := &ListenResponse{}

	response.Events = listener.Listen()

	res = response

	return
}

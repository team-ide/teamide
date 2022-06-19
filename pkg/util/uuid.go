package util

import (
	"github.com/google/uuid"
	"strings"
)

// UUID 生成UUID
func UUID() (res string) {
	res = uuid.NewString()
	res = strings.ReplaceAll(res, "-", "")
	return
}

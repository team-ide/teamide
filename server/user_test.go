package main

import (
	"server/service/userService"
	"testing"
)

func TestUser(t *testing.T) {
	userService.TestTotalBatchInsert(30, 200000)
}

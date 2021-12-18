package main

import (
	userService "server/service/user"
	"testing"
)

func TestUser(t *testing.T) {
	userService.TestTotalBatchInsert(30, 200000)
}

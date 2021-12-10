package main

import (
	"server/userService"
	"testing"
)

func TestUser(t *testing.T) {
	userService.TestTotalBatchInsert(30, 200000)
}

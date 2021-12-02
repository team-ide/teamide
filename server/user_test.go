package main

import (
	"service"
	"testing"
)

func TestUser(t *testing.T) {
	service.TestTotalBatchInsert(30, 200000)
}

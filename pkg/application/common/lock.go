package common

type Locker interface {
	Lock() error
	Unlock() error
}

package core

import (
	"github.com/xakpro/sms/internal/interfaces"
)

type Core struct {
	lg               interfaces.Logger
	cache            interfaces.Cache
	smscUsername     string
	smscPassword     string
	smscSender       string
	balanceNotifUrls map[float64]string
}

func NewCore(lg interfaces.Logger, cache interfaces.Cache, smscUsername string, smscPassword string, smscSender string, array map[float64]string) *Core {
	core := &Core{
		lg:               lg,
		cache:            cache,
		smscUsername:     smscUsername,
		smscPassword:     smscPassword,
		smscSender:       smscSender,
		balanceNotifUrls: array,
	}

	return core
}

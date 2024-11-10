package server

import (
	"fmt"
	"time"
)

type cache struct {
	values map[string]string
}

var c cache

func InitCache() {
	c = cache{
		values: make(map[string]string),
	}
}

func setValue(cmd setCommand) {
	if cmd.ttl > 0 {
		timer := time.NewTimer(cmd.ttl)
		go func(key string, timer *time.Timer) {
			select {
			case <-timer.C:
				delete(c.values, key)
				LogInfo(fmt.Sprintf("%s expired", key))
			}
		}(cmd.key, timer)
	}
	c.values[cmd.key] = cmd.value
	LogInfo(fmt.Sprintf("Saved %s with ttl of %v to cache", cmd.key, cmd.ttl))
}

func getValue(key string) (string, bool) {
	val, ok := c.values[key]
	return val, ok
}

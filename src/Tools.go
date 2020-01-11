package src

import (
	guuid "github.com/satori/go.uuid"
	"time"
)

func DoEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

func GetHash() string {
	id, _ := guuid.NewV4()
	return id.String()
}

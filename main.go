package main

import (
	"fmt"
	"sync/atomic"

	"github.com/MarinX/keylogger"
	"github.com/RamadanIbrahem98/kabsa/keys"
	"github.com/robfig/cron/v3"
)

type Kabsa struct {
	times int64
}

func persist(kabsa *Kabsa) {
	letters := atomic.LoadInt64(&kabsa.times)
	atomic.StoreInt64(&kabsa.times, 0)

	fmt.Println("Persisting ", letters, " presses, so the average WPM is ", letters/5)
}

func main() {
	keyboard := keylogger.FindKeyboardDevice()

	if len(keyboard) <= 0 {
		panic("No keyboard found")
	}

	fmt.Println("Found a keyboard at", keyboard)
	k, err := keylogger.New(keyboard)
	if err != nil {
		panic(err)
	}
	defer k.Close()

	kabsa := &Kabsa{times: 0}

	c := cron.New()
	c.AddFunc("* * * * *", func() {
		go persist(kabsa)
	}) // run every minute
	c.Start()
	defer c.Stop()

	events := k.Read()

	for e := range events {
		if e.Type == keylogger.EvKey && e.KeyRelease() {
			keyString := keys.KeyCodeMap[e.Code]
			if keys.CountableKeys[keyString] {
				atomic.AddInt64(&kabsa.times, 1)
				fmt.Println("[event] release key ", e.KeyString(), " you have pressed ", kabsa.times, " presses.")
			}
		}
	}
}

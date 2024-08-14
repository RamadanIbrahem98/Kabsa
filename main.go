package main

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/MarinX/keylogger"
	"github.com/RamadanIbrahem98/kabsa/keys"
)

type Kabsa struct {
	times   int64
	startAt int64
}

var watchdogTimer *time.Timer

func resetWatchdog(kabsa *Kabsa) {
	if watchdogTimer != nil {
		watchdogTimer.Stop()
	}
	watchdogTimer = time.AfterFunc(1*time.Second, func() {
		endAt := time.Now().UnixMilli()
		startAt := atomic.LoadInt64(&kabsa.startAt)

		if startAt == 0 {
			return
		}
		atomic.StoreInt64(&kabsa.startAt, 0)
		letters := atomic.LoadInt64(&kabsa.times)
		atomic.StoreInt64(&kabsa.times, 0)
		elapsedTime := (endAt - startAt) / 1000.0
		wpm := int64((float64(letters) / 5.0) / (float64(elapsedTime) / 60.0))

		fmt.Println("Watchdog triggered, so the average WPM is ", wpm)
	})
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

	events := k.Read()

	for e := range events {
		if e.Type == keylogger.EvKey && e.KeyRelease() {
			keyString := keys.KeyCodeMap[e.Code]
			if keys.CountableKeys[keyString] {
				atomic.AddInt64(&kabsa.times, 1)
				startAt := atomic.LoadInt64(&kabsa.startAt)

				if startAt == 0 {
					atomic.StoreInt64(&kabsa.startAt, time.Now().UnixMilli())
				}

				resetWatchdog(kabsa)
			}
		}
	}
}

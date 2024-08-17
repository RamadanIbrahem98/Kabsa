package main

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/MarinX/keylogger"
	"github.com/RamadanIbrahem98/kabsa/db"
	"github.com/RamadanIbrahem98/kabsa/keys"
)

type Kabsa struct {
	times   int64
	startAt int64
	endAt   int64
}

var debouncer *time.Timer

func resetWatchdog(kabsa *Kabsa, db *db.DB) {
	if debouncer != nil {
		debouncer.Stop()
	}
	debouncer = time.AfterFunc(5*time.Second, func() {
		endAt := atomic.LoadInt64(&kabsa.endAt) + 1000
		startAt := atomic.LoadInt64(&kabsa.startAt)

		if startAt == 0 || endAt == 1000 {
			return
		}

		atomic.StoreInt64(&kabsa.startAt, 0)
		atomic.StoreInt64(&kabsa.endAt, 0)

		letters := atomic.LoadInt64(&kabsa.times)
		atomic.StoreInt64(&kabsa.times, 0)

		elapsedTime := (endAt - startAt) / 1000.0
		wpm := int64((float64(letters) / 5.0) / (float64(elapsedTime) / 60.0))

		db.Insert(letters, startAt, endAt, wpm)

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

	db, err := db.New()

	if err != nil {
		panic(err)
	}

	defer db.Close()

	events := k.Read()

	for e := range events {
		if e.Type == keylogger.EvKey && e.KeyRelease() {
			keyString := keys.KeyCodeMap[e.Code]
			if keys.CountableKeys[keyString] {
				atomic.AddInt64(&kabsa.times, 1)
				atomic.StoreInt64(&kabsa.endAt, time.Now().UnixMilli())
				startAt := atomic.LoadInt64(&kabsa.startAt)

				if startAt == 0 {
					atomic.StoreInt64(&kabsa.startAt, time.Now().UnixMilli())
				}

				resetWatchdog(kabsa, db)
			}
		}
	}
}

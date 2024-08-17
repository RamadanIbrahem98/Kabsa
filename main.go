package main

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/MarinX/keylogger"
	"github.com/RamadanIbrahem98/kabsa/kabsa"
	"github.com/RamadanIbrahem98/kabsa/keys"
)

var debouncer *time.Timer

func resetWatchdog(kabsa *kabsa.Kabsa) {
	if debouncer != nil {
		debouncer.Stop()
	}
	debouncer = time.AfterFunc(5*time.Second, func() {
		endAt := atomic.LoadInt64(&kabsa.EndAt) + 1000
		startAt := atomic.LoadInt64(&kabsa.StartAt)

		if startAt == 0 || endAt == 1000 {
			return
		}

		atomic.StoreInt64(&kabsa.StartAt, 0)
		atomic.StoreInt64(&kabsa.EndAt, 0)

		letters := atomic.LoadInt64(&kabsa.Presses)
		atomic.StoreInt64(&kabsa.Presses, 0)

		elapsedTime := (endAt - startAt) / 1000.0
		wpm := int64((float64(letters) / 5.0) / (float64(elapsedTime) / 60.0))

		kabsa.DB.Insert(letters, startAt, endAt, wpm)

		fmt.Println("Debouncer timedout, so the average WPM is ", wpm)
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

	kabsa, err := kabsa.New()

	if err != nil {
		panic(err)
	}

	defer kabsa.DB.Close()

	events := k.Read()

	for e := range events {
		if e.Type == keylogger.EvKey && e.KeyRelease() {
			keyString := keys.KeyCodeMap[e.Code]
			if keys.CountableKeys[keyString] {
				atomic.AddInt64(&kabsa.Presses, 1)
				atomic.StoreInt64(&kabsa.EndAt, time.Now().UnixMilli())
				startAt := atomic.LoadInt64(&kabsa.StartAt)

				if startAt == 0 {
					atomic.StoreInt64(&kabsa.StartAt, time.Now().UnixMilli())
				}

				resetWatchdog(kabsa)
			}
		}
	}
}

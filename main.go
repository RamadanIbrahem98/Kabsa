package main

import (
	"fmt"
	"sync/atomic"

	"github.com/MarinX/keylogger"
	"github.com/RamadanIbrahem98/kabsa/keys"
)

type Kabsa struct {
	times int64
}

func main() {
	keyboard := keylogger.FindKeyboardDevice()

	if len(keyboard) <= 0 {
		panic("No keyboard found...you will need to provide manual input path")
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
				fmt.Println("[event] release key ", e.KeyString(), " you have pressed ", kabsa.times, " presses.")
			}
		}
	}
}

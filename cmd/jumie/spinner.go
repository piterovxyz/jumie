package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	colorGreen = "\033[32m"
	colorReset = "\033[0m"
)

var blinkFrames = []string{
	colorGreen + "(o_o)" + colorReset,
	colorGreen + "(-_-)" + colorReset,
	colorGreen + "(o_o)" + colorReset,
	colorGreen + "(-_-)" + colorReset,
	colorGreen + "(o_o)" + colorReset,
	colorGreen + "(-_-)" + colorReset,
	colorGreen + "(o_o)" + colorReset,
}

var thinkFrames = []string{
	colorGreen + "(o_o)" + colorReset,
	colorGreen + "(o_O)" + colorReset,
	colorGreen + "(O_o)" + colorReset,
}

var tips = []string{
	"thinking...",
	"pondering...",
	"contemplating...",
	"analyzing...",
	"processing...",
	"scheming...",
	"brainstorming...",
	"deliberating...",
	"evaluating...",
	"reflecting...",
	"plotting...",
	"wondering...",
	"calculating...",
	"reasoning...",
	"interpreting...",
	"visualizing...",
	"deciphering...",
	"figuring out...",
	"synthesizing...",
	"philosophizing...",
}

func startSpinner() func() {
	stop := make(chan struct{})
	done := make(chan struct{})

	go func() {
		defer close(done)
		tip := tips[rand.Intn(len(tips))]
		runes := []rune(tip)
		tipLen := len(runes)

		ticker := time.NewTicker(40 * time.Millisecond)
		defer ticker.Stop()

		step := 0
		blinkTicks := len(blinkFrames) * 5

		for {
			select {
			case <-stop:
				fmt.Print("\r\033[K")
				return
			case <-ticker.C:
				var currentFrame string
				if step < blinkTicks {
					currentFrame = blinkFrames[step/5]
				} else {
					thinkTick := (step - blinkTicks) / 4
					currentFrame = thinkFrames[thinkTick%len(thinkFrames)]
				}

				charsToShow := step
				if charsToShow > tipLen {
					charsToShow = tipLen
				}
				currentText := string(runes[:charsToShow])

				fmt.Printf("\r\033[K%s %s", currentFrame, currentText)
				step++
			}
		}
	}()

	return func() {
		close(stop)
		<-done
	}
}

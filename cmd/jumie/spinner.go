package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Dim       = "\033[2m"
	Gray      = "\033[90m"
	Cyan      = "\033[36m"
	Green     = "\033[32m"
	Yellow    = "\033[33m"
	BoldGreen = "\033[1;32m"
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

func startPromptSpinner() func() {
	stop := make(chan struct{})
	done := make(chan struct{})

	fmt.Printf("\n%s(o_o) %sexecute? (y/n): %s\0337", colorGreen, BoldGreen, Reset)

	go func() {
		defer close(done)
		for {
			select {
			case <-stop:
				return
			case <-time.After(1 * time.Second):
				fmt.Printf("\0337\r%s(-_-)%s\0338", colorGreen, Reset)

				select {
				case <-stop:
					return
				case <-time.After(150 * time.Millisecond):
					fmt.Printf("\0337\r%s(o_o)%s\0338", colorGreen, Reset)
				}
			}
		}
	}()

	return func() {
		close(stop)
		<-done
	}
}

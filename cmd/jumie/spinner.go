package main

import (
	"fmt"
	"time"
)

const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Cyan      = "\033[36m"
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

func startSpinner() (func(), func(string)) {
	stopChan := make(chan struct{})
	done := make(chan struct{})
	tipChan := make(chan string, 1)

	go func() {
		defer close(done)

		tip := "..."

		runes := []rune(tip)
		tipLen := len(runes)

		ticker := time.NewTicker(40 * time.Millisecond)
		defer ticker.Stop()

		step := 0
		blinkTicks := len(blinkFrames) * 5

		for {
			select {
			case <-stopChan:
				fmt.Print("\r\033[K")
				return
			case newTip := <-tipChan:
				tip = newTip
				runes = []rune(tip)
				tipLen = len(runes)
				step = 0
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

	stopFunc := func() {
		close(stopChan)
		<-done
	}

	updateFunc := func(newTip string) {
		select {
		case tipChan <- newTip:
		default:
		}
	}

	return stopFunc, updateFunc
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

func startLoginPrompt() func() {
	stop := make(chan struct{})
	done := make(chan struct{})

	fmt.Printf("%s(o_o) %senter your gemini api key: %s\0337", colorGreen, BoldGreen, Reset)

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

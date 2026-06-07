package main

import (
	"bufio"
	"fmt"
	"jumie/internal/ai"
	"jumie/internal/ipc"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: %s <message> or %s login <key>", os.Args[0], os.Args[0])
	}

	msg := strings.Join(os.Args[1:], " ")

	c, err := ipc.NewClient()
	if err != nil {
		log.Fatalf("error creating ipc client: %v\n", err)
	}

	stop := startSpinner()
	resp, err := c.RequestPlan(msg)
	stop()

	if err != nil {
		log.Fatalf("error sending message: %v\n", err)
	}

	confirm := do(resp)

	if !confirm {
		os.Exit(0)
	}

	err = c.DoPlan(resp)
	if err != nil {
		log.Fatalf("error executing plan: %v\n", err)
	}
}

func typewriterPrint(text string, delay time.Duration) {
	for _, r := range []rune(text) {
		fmt.Print(string(r))
		time.Sleep(delay)
	}
	fmt.Println()
}

func do(plan *ai.Plan) bool {
	if len(plan.Steps) == 0 {
		fmt.Println(Yellow + "plan is empty" + Reset)
		return false
	}

	fmt.Println()
	fmt.Println(Bold + Cyan + "✦ jumie plan:" + Reset)

	for _, step := range plan.Steps {
		fmt.Print(Cyan + "➜  " + Reset)
		typewriterPrint(step.Description, 6*time.Millisecond)

		fmt.Printf("%s$  %s%s\n", Yellow, step.Command, Reset)
	}

	reader := bufio.NewReader(os.Stdin)

	stopPrompt := startPromptSpinner()
	input, err := reader.ReadString('\n')
	stopPrompt()

	if err != nil {
		return false
	}

	input = strings.TrimSpace(strings.ToLower(input))
	return input == "y" || input == "yes"
}

package main

import (
	"bufio"
	"context"
	"fmt"
	"jumie/internal/ai"
	"jumie/internal/config"
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

	if os.Args[1] == "login" {
		if len(os.Args) < 3 {
			log.Fatalf("usage: %s login <api_key>", os.Args[0])
		}
		key := os.Args[2]

		fmt.Println("validating API key...")
		client, err := ai.NewClient(key)
		if err != nil {
			log.Fatalf("failed to initialize AI client: %v\n", err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := client.ValidateKey(ctx); err != nil {
			cancel()
			log.Fatalf("invalid API key: %v\n", err)
		}
		cancel()

		if err := config.Save(key); err != nil {
			log.Fatalf("error saving config: %v\n", err)
		}
		fmt.Println("successfully logged in!")
		return
	}

	msg := strings.Join(os.Args[1:], " ")

	cfg, err := config.Load()
	if err != nil || cfg.APIKey == "" {
		stopPrompt := startLoginPrompt()
		reader := bufio.NewReader(os.Stdin)
		key, err := reader.ReadString('\n')
		stopPrompt()
		if err != nil {
			log.Fatalf("\nerror reading api key: %v\n", err)
		}
		key = strings.TrimSpace(key)
		if key == "" {
			log.Fatalf("\napi key cannot be empty\n")
		}

		fmt.Println("validating API key...")
		client, err := ai.NewClient(key)
		if err != nil {
			log.Fatalf("failed to initialize AI client: %v\n", err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := client.ValidateKey(ctx); err != nil {
			cancel()
			fmt.Printf("invalid api key! please provide a valid api key or try again later.\n")
			return
		}
		cancel()

		if err := config.Save(key); err != nil {
			log.Fatalf("\nerror saving config: %v\n", err)
		}
		fmt.Println("successfully logged in!")
	}

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

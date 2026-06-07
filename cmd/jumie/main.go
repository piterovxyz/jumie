package main

import (
	"bufio"
	"fmt"
	"jumie/internal/ai"
	"jumie/internal/daemon"
	"jumie/internal/installer"
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

	if err := c.Ping(); err != nil {
		log.Fatalf("daemon is not running! please start jumied first.\n")
	}

	if err := checkDeps(c); err != nil {
		log.Fatalf("dependency setup aborted: %v\n", err)
	}

	stop, updateTip := startSpinner()
	resp, err := c.RequestPlan(msg, updateTip)
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
	if plan.Reasoning != "" {
		fmt.Print("\033[2m")
		typewriterPrint("✦ jumie reasoning: "+plan.Reasoning, 3*time.Millisecond)
		fmt.Print("\033[0m\n")
	}

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

func checkDeps(c *ipc.Client) error {
	err := c.StartOllama()
	if err == nil {
		return nil
	}

	if !strings.Contains(err.Error(), "not installed") {
		return fmt.Errorf("daemon ollama error: %v", err)
	}

	fmt.Print(Yellow + "jumie is required to install ollama. install now? (y/n): " + Reset)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	input = strings.TrimSpace(strings.ToLower(input))
	if input != "y" && input != "yes" {
		return fmt.Errorf("installation aborted")
	}

	err = installer.InstallOllama(func(p string) {
		fmt.Println(Cyan + p + Reset)
	})
	if err != nil {
		return err
	}

	err = c.StartOllama()
	if err != nil {
		return fmt.Errorf("failed to start ollama after install: %v", err)
	}

	err = daemon.PullModel(func(p string) {
		fmt.Print(Cyan + p + Reset)
	})
	if err != nil {
		return err
	}

	return nil
}

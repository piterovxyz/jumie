package indexer

import (
	"bufio"
	"bytes"
	"errors"
	"log"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
)

func RunIndexer(c *InfoCache) {
	osType, osRelease, err := getSystemRelease()
	if err != nil {
		return
	}

	current, err := user.Current()
	if err != nil {
		log.Printf("error to get current user: %v\n", err)
		return
	}

	isRoot := current.Uid == "0"

	info := SystemInfo{
		osType,
		osRelease,
		os.Getenv("SHELL"),
		isRoot,
		make(map[string]string),
	}

	c.Write(info)
}

func CheckBinaries(tools []string) map[string]bool {
	res := make(map[string]bool)
	for _, tool := range tools {
		_, err := exec.LookPath(tool)
		res[tool] = err == nil
	}
	return res
}

func getSystemRelease() (string, string, error) {
	var out bytes.Buffer
	switch runtime.GOOS {
	case "darwin":
		cmd := exec.Command("sw_vers", "-productVersion")
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			log.Printf("error to get system version: %v\n", err)
			return "darwin", "unsupported", err
		}
		return "darwin", strings.TrimSpace(out.String()), nil
	case "linux":
		file, err := os.Open("/etc/os-release")
		if err != nil {
			log.Printf("error to get system version: %v\n", err)
			return "linux", "unsupported", err
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Printf("error to get system version: %v\n", err)
			}
		}(file)

		sc := bufio.NewScanner(file)
		for sc.Scan() {
			line := sc.Text()
			if version, ok := strings.CutPrefix(line, "VERSION_ID="); ok {
				return "linux", strings.Trim(version, "\""), nil
			}
		}

		if sc.Err() != nil {
			log.Printf("error to get system version: %v\n", sc.Err())
			return "linux", "unsupported", err
		}
		return "linux", "unsupported", err
	default:
		return "unknown", "unsupported", errors.New("unsupported system")
	}
}

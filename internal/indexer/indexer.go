package indexer

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func RunIndexer() {
	osType, osRelease := getSystemRelease()
	info := SystemInfo{
		osType,
		osRelease,
		parsePath(),
		make(map[string]string),
	}

	fmt.Println(info)
}

func getSystemRelease() (string, string) {
	var out bytes.Buffer
	switch runtime.GOOS {
	case "darwin":
		cmd := exec.Command("sw_vers", "-productVersion")
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			log.Printf("error to get system version: %v\n", err)
			return "darwin", "unsupported"
		}
		return "darwin", strings.TrimSpace(out.String())
	case "linux":
		file, err := os.Open("/etc/os-release")
		if err != nil {
			log.Printf("error to get system version: %v\n", err)
			return "linux", "unsupported"
		}
		defer file.Close()

		sc := bufio.NewScanner(file)
		for sc.Scan() {
			if sc.Err() != nil {
				log.Printf("error to get system version: %v\n", err)
				return "linux", "unsupported"
			}

			line := sc.Text()
			if version, ok := strings.CutPrefix(line, "VERSION_ID="); ok {
				return "linux", version
			}
		}
		return "linux", "unsupported"
	default:
		return "unknown", "unsupported"
	}
}

func parsePath() []string {
	var res []string
	path := os.Getenv("PATH")
	folders := strings.SplitSeq(path, ":")

	for f := range folders {
		bins, err := os.ReadDir(f)
		if err != nil {
			log.Printf("error reading folder: %v\n", err)
			continue
		}

		for _, bin := range bins {
			if bin.IsDir() {
				continue
			}

			info, err := bin.Info()
			if err != nil {
				continue
			}

			if info.Mode()&0111 != 0 {
				res = append(res, bin.Name())
			}
		}
	}

	return res
}

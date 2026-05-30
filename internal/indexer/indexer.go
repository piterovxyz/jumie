package indexer

import (
	"bufio"
	"bytes"
	"errors"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func RunIndexer(c *InfoCache) {
	osType, osRelease, err := getSystemRelease()
	if err != nil {
		return
	}

	path, err := parsePath()
	if err != nil {
		return
	}

	info := SystemInfo{
		osType,
		osRelease,
		path,
		os.Getenv("SHELL"),
		make(map[string]string),
	}

	c.Write(info)
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
		defer file.Close()

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

func parsePath() ([]string, error) {
	var res []string
	path := os.Getenv("PATH")
	folders := strings.SplitSeq(path, ":")

	for f := range folders {
		bins, err := os.ReadDir(f)
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				log.Printf("error reading folder: %v\n", err)
			}
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

	return res, nil
}

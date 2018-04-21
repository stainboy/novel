package util

import (
	"bufio"
	"os"
	"strings"
)

func Source(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Index(line, "#") == 0 {
			continue
		}
		kv := strings.SplitN(line, "=", 2)
		if len(kv) == 2 {
			os.Setenv(kv[0], kv[1])
		}
	}

	return nil
}

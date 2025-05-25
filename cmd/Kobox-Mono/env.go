package main

import (
	"fmt"
	"os"
	"strings"
)

func updateEnvFile(key, value string) error {
	content, err := os.ReadFile(envFile)
	if err != nil && os.IsNotExist(err) {
		return err
	}
	updated := false
	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, key) {
			lines[i] = fmt.Sprintf("%s=%s", key, value)
			updated = true
			break
		}
	}
	if !updated {
		lines = append(lines, fmt.Sprintf("%s=%s", key, value))
	}

	return os.WriteFile(envFile, []byte(strings.Join(lines, "\n")), 0644)
}

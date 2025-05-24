package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func ConvertFileWithPandoc(inputFile, outputFile, fromFormat, toFormat string) error {
	cmd := exec.Command("pandoc",
		"-f", fromFormat,
		"-t", toFormat,
		"-o", outputFile,
		inputFile)

	return cmd.Run()
}

func ConvertStringWithPandoc(content, fromFormat, toFormat string) ([]byte, error) {
	cmd := exec.Command("pandoc", "-f", fromFormat, "-t", toFormat)

	cmd.Stdin = strings.NewReader(content)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return []byte{}, err
	}

	return output, nil
}

func GenerateHTML(title, body string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s</title>
</head>
<body>
    %s
</body>
</html>`, title, body)
}

func GetAppDataDir(appname string) (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	appDir := filepath.Join(configDir, appname)
	err = os.Mkdir(appDir, 0755)
	if err != nil {
		return "", err
	}

	return appDir, nil
}

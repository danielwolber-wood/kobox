package main

import (
	"fmt"
	"os/exec"
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

func ConvertStringWithPandoc(content HTML, fromFormat, toFormat string) (Epub, error) {
	cmd := exec.Command("pandoc", "-f", fromFormat, "-t", toFormat)

	cmd.Stdin = strings.NewReader(string(content))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return []byte{}, err
	}

	return output, nil
}

func GenerateHTML(title, body string) HTML {
	return HTML(fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s</title>
</head>
<body>
    %s
</body>
</html>`, title, body))
}

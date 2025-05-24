package main

import (
	"github.com/mmcdole/gofeed"
	"os"
	"path/filepath"
)

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

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func ParseWithGofeed(url string) (*gofeed.Feed, error) {
	parser := gofeed.NewParser()
	feed, err := parser.ParseURL(url)
	if err != nil {
		return nil, err
	}
	return feed, nil
}

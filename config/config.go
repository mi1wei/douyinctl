package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var configPath = "./config.json"

type SearchParams struct {
	Keywords        []string
	ScrollPageTimes int
	SourceFilePaths []string
}

type CommentParams struct {
	Keywords        []string
	ScrollPageTimes int
	SourceFilePaths []string
}

type Account struct {
	UserId         string
	UserName       string
	Count          int32
	CookieFilePath string
	SearchParams   SearchParams
	CommentParams  CommentParams
}

type Config struct {
	Accounts []*Account
}

func (c *Config) Store() error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(configPath, data, 0644)
	if err != nil {
		return err
	}

	fmt.Println("Config stored to file:", configPath)
	return nil
}

func (c *Config) Load() error {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, c)
	if err != nil {
		return err
	}

	fmt.Println("Config loaded from file:", configPath)
	return nil
}

var DefaultConfig = &Config{
	Accounts: []*Account{
		{
			UserId: "aa",
			Count:  0,
			SearchParams: SearchParams{
				Keywords:        []string{},
				ScrollPageTimes: 5,
				SourceFilePaths: []string{},
			},
			CommentParams: CommentParams{
				Keywords:        []string{},
				ScrollPageTimes: 5,
				SourceFilePaths: []string{},
			},
		},
		{
			UserId: "bb",
			Count:  0,
			SearchParams: SearchParams{
				Keywords:        []string{},
				ScrollPageTimes: 5,
				SourceFilePaths: []string{},
			},
			CommentParams: CommentParams{
				Keywords:        []string{},
				ScrollPageTimes: 5,
				SourceFilePaths: []string{},
			},
		},
	},
}

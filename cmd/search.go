package cmd

import (
	"fmt"
	"github.com/mi1wei/douyinctl/config"
	"github.com/spf13/cobra"
	"log"
	"path/filepath"
	"sync"
)

const searchUrl = "https://www.douyin.com/search/%s?publish_time=182&sort_type=0&source=tab_search&type=video"

type SearchOption struct {
	config *config.Config
}

var keywords []string
var scrollPage int
var sourceFiles []string

var searchCmd = &cobra.Command{
	Use:     "search",
	Short:   "Search link from DouYin,XiaoHongShu",
	Example: "douyinctl search --keyword keyword1,keyword2",
	Run: func(cmd *cobra.Command, args []string) {
		err := search()
		if err != nil {
			fmt.Println("Error: search failed, err: ", err)
			return
		}
		fmt.Println("Search successful")
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().StringSliceVar(&keywords, "keyword", []string{}, "A list of keywords to search, such as 广东,夜市")
	searchCmd.Flags().IntVar(&scrollPage, "scroll", 0, "Scroll page times for each keyword")
	searchCmd.Flags().StringSliceVar(&sourceFiles, "sourceFile", []string{}, "The source file for search")
}

func search() error {
	opts := &SearchOption{
		config: &config.Config{},
	}
	if err := opts.config.Load(); err != nil {
		config.DefaultConfig.Store()
		return err
	}

	for _, account := range opts.config.Accounts {
		if len(keywords) > 0 {
			account.SearchParams.Keywords = keywords
		}

		if scrollPage > 0 {
			account.SearchParams.ScrollPageTimes = scrollPage
		}

		if len(sourceFiles) > 0 {
			account.SearchParams.SourceFilePaths = sourceFiles
		}
	}

	results, _ := opts.performSearch()
	for account, links := range results {
		fmt.Println(account.UserName, len(links))
	}

	return nil
}

func (opts *SearchOption) performSearch() (map[*config.Account][]string, error) {
	searchResutls := make(map[*config.Account][]string, 0)
	for _, account := range opts.config.Accounts {
		mutex := sync.Mutex{}
		links := make([]string, 0)
		var wg sync.WaitGroup
		errors := make(map[string]error)

		for _, keyword := range account.SearchParams.Keywords {
			wg.Add(1)
			go func(keyword string) {
				defer wg.Done()

				items, err := opts.ListItems(&ListItemsParams{
					Keyword: keyword,
					Url:     fmt.Sprintf(searchUrl, keyword),
				}, account)
				if err != nil {
					errors[keyword] = err
					return
				}

				mutex.Lock()
				links = append(links, items...)
				mutex.Unlock()

			}(keyword)
		}

		for _, sourceFile := range sourceFiles {
			wg.Add(1)
			go func(sourceFile string) {
				prefix := filepath.Base(sourceFile)
				defer wg.Done()

				items, err := opts.ListItems(&ListItemsParams{
					sourceFile: sourceFile,
				}, account)
				if err != nil {
					errors[prefix] = err
					return
				}

				mutex.Lock()
				links = append(links, items...)
				mutex.Unlock()

			}(sourceFile)
		}
		wg.Wait()

		if len(errors) > 0 {
			for keyword, err := range errors {
				log.Printf("Error for %s: %v\n", keyword, err)
			}
		}

		searchResutls[account] = links

	}

	return searchResutls, nil
}

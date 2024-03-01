package cmd

import (
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/mi1wei/douyinctl/config"
	"github.com/mi1wei/douyinctl/util/cookie"
	"github.com/mi1wei/douyinctl/util/file"
	"github.com/mi1wei/douyinctl/util/my_chromedp"
	"strings"
)

const (
	videoSelector = "#search-content-area > div > div.HHwqaG_P > div:nth-child(2) > ul > li > div > a"
)

type ListItemsParams struct {
	Keyword         string
	Url             string
	ScrollPageTimes int
	sourceFile      string
}

func (opts *SearchOption) ListItems(params *ListItemsParams, account *config.Account) ([]string, error) {
	ctx, cancel := my_chromedp.NewContext(false)
	defer cancel()

	//loginStatusSelector := "#douyin-header > div.oJArD0aS > header > div > div > div.iqHX00br > div > div > div:nth-child(5) > div > a > div > img"

	if len(params.sourceFile) > 0 {
		hrefs, err := file.Read(params.sourceFile)
		if err != nil {
			return nil, err
		}
		return hrefs, nil
	}

	if err := chromedp.Run(ctx, cookie.LoadCookies("douyin", account.UserId), chromedp.Navigate(params.Url)); err != nil {
		return nil, err
	}

	my_chromedp.Scroll(ctx, account.SearchParams.ScrollPageTimes)

	hrefs, _, err := my_chromedp.Texts(ctx, videoSelector, func(nodes []*cdp.Node) []string {
		hrefs := make([]string, 0)
		for _, node := range nodes {
			if strings.Contains(node.Attributes[1], "note") {
				fmt.Println("note")
				continue
			}
			hrefs = append(hrefs, fmt.Sprintf("https:%s", node.Attributes[1]))
		}
		return hrefs
	})
	if err != nil {
		return nil, err
	}

	return hrefs, nil
}

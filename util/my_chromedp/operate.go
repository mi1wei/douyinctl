package my_chromedp

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/mi1wei/douyinctl/constants"
)

type WithNodeOperation = func([]*cdp.Node) []string

func Scroll(ctx context.Context, times int) error {
	scrollResp := ""
	for i := 0; i < times; i++ {
		if times > 10 && i%10 == 0 {
			fmt.Println("refreshing %d page...", i)
		}
		chromedp.Run(ctx, chromedp.Sleep(constants.DefaultScrollTimeout), chromedp.Evaluate(constants.ScrollScript, &scrollResp))
	}
	return nil
}

func Nodes(ctx context.Context, selector string) ([]*cdp.Node, error) {
	ctx, cancel := context.WithTimeout(ctx, constants.Timeout)
	defer cancel()
	nodes := make([]*cdp.Node, 0)
	if err := chromedp.Run(ctx, chromedp.Nodes(selector, &nodes, chromedp.ByQueryAll)); err != nil {
		return nodes, err
	}
	return nodes, nil
}

func Texts(ctx context.Context, selector string, operation WithNodeOperation) ([]string, int, error) {
	texts := make([]string, 0)
	nodes, err := Nodes(ctx, selector)
	if err != nil {
		return texts, 0, err
	}

	if operation != nil {
		texts = operation(nodes)
	}

	return texts, len(texts), nil
}

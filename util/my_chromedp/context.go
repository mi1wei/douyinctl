package my_chromedp

import (
	"context"
	"github.com/chromedp/chromedp"
	"log"
)

func NewContext(flag bool) (context.Context, context.CancelFunc) {
	allocatorOptions := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", flag),
		//chromedp.WindowSize(100, 200),
		chromedp.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36"), // 设置User-Agent
	}
	options := append(chromedp.DefaultExecAllocatorOptions[:], allocatorOptions...)
	execAllocatorCtx, _ := chromedp.NewExecAllocator(context.Background(), options...)
	ctx, cancel := chromedp.NewContext(execAllocatorCtx, chromedp.WithLogf(log.Printf))
	return ctx, cancel
}

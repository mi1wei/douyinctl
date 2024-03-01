package constants

import "time"

const (
	ScrollScript = `window.scrollTo(0, document.body.scrollHeight);`
	Timeout      = 10 * time.Second
	Concurrency  = 2

	PageCountSelector = "#search-content-area > div > div.HHwqaG_P > div:nth-child(2) > ul > li > div > div > a > div > div.F84uEzbg > div > div.swoZuiEM"
	DataStorePath     = "/Users/miwei/Desktop/data/"
	DevToolsWsURL     = "ws://127.0.0.1:9222/devtools/browser/28b1d22e-6d04-4be0-b81d-d388c9934da6"

	DefaultScrollTimeout = 1 * time.Second
)

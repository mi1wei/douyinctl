package cmd

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"github.com/mi1wei/douyinctl/config"
	"github.com/mi1wei/douyinctl/util/chatgpt"
	"github.com/mi1wei/douyinctl/util/cookie"
	"github.com/mi1wei/douyinctl/util/my_chromedp"
	"github.com/mi1wei/douyinctl/util/my_string"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"time"
)

type CommentOption struct {
	ctx    context.Context
	cancel context.CancelFunc
	config *config.Config
}

var commentCmd = &cobra.Command{
	Use:     "comment",
	Short:   "Login to platform",
	Example: "douyinctl login --platform douyin",
	Run: func(cmd *cobra.Command, args []string) {
		err := comment()
		if err != nil {
			fmt.Println("Error: comment failed, err: ", err)
			return
		}
		fmt.Println("Comment successful")
	},
}

func init() {
	rootCmd.AddCommand(commentCmd)
	commentCmd.Flags().IntVar(&scrollPage, "scroll", 0, "Scroll page times for each keyword")
	commentCmd.Flags().StringSliceVar(&sourceFiles, "sourceFile", []string{}, "The source file for search")
}

func comment() error {
	opts := &CommentOption{
		config: &config.Config{},
	}

	if err := opts.config.Load(); err != nil {
		config.DefaultConfig.Store()
		return nil
	}

	for _, account := range opts.config.Accounts {
		if scrollPage > 0 {
			account.CommentParams.ScrollPageTimes = scrollPage
		}

		if len(sourceFiles) > 0 {
			account.CommentParams.SourceFilePaths = sourceFiles
		}
	}

	addCommentesParams := AddCommentesParams{}

	if len(sourceFiles) == 0 {
		searchResults, _ := (&SearchOption{config: opts.config}).performSearch()
		addCommentesParams.AccountVideoUrls = searchResults
	}

	return opts.addComments(addCommentesParams)
}

type AddCommentesParams struct {
	AccountVideoUrls map[*config.Account][]string
	videoUrls        []string
}

func (opts *CommentOption) addComments(params AddCommentesParams) error {
	var (
		nodesAll  []*cdp.Node
		gptClient = chatgpt.NewClient()
		selector  = "#douyin-right-container > div:nth-child(2) > div > div.leftContainer.gkVJg5wr > div.KwRNeXA3 > div > div > div.HV3aiR5J.comment-mainContent > div > div > div.VjrdhTqP > div > div.LvAtyU_f > span > span > span > span > span > span > span"
	)

	for _, account := range opts.config.Accounts {
		var urls []string

		if len(params.AccountVideoUrls) > 0 {
			urls = params.AccountVideoUrls[account]
		} else {
			urls = params.videoUrls
		}

		for _, videoUrl := range urls {
			ctx, cancel := my_chromedp.NewContext(false)
			defer cancel()

			logging := logrus.WithField("account", account)

			if err := chromedp.Run(ctx, cookie.LoadCookies("douyin", account.UserId), chromedp.Navigate(videoUrl)); err != nil {
				return err
			}

			if err := chromedp.Run(ctx, chromedp.Nodes(selector, &nodesAll, chromedp.ByQueryAll)); err != nil {
				return err
			}

			var title string
			var titleSelector = "#douyin-right-container > div:nth-child(2) > div > div.leftContainer.gkVJg5wr > div.XYnWH9QO > div > div.x_vgJ3yL.kDzTQY11 > div > h1 > span > span:nth-child(2) > span > span:nth-child(1) > span > span > span"
			if err := chromedp.Run(ctx, chromedp.Text(titleSelector, &title)); err != nil {
				return err
			}

			logging.Infof("[title: %s]", title)
			if firstPromot, err := gptClient.FirstMessage(title); err != nil {
				logging.Errorf("[first promot failed : %v]", err)
			} else {
				logging.Infof("[first promot resp: %s]", firstPromot)
			}

			for i, node := range nodesAll {
				if len(node.Children) <= 0 {
					continue
				}

				if len(account.CommentParams.Keywords) > 0 {
					if !my_string.Contains(account.CommentParams.Keywords, node.Children[0].NodeValue) {
						continue
					}
				} else if i > 10 {
					break
				}

				time.Sleep(30 * time.Second)
				reply, err := gptClient.Send(fmt.Sprintf(chatgpt.REQUEST_PREFIX, node.Children[0].NodeValue))
				if err != nil {
				}

				logging.Infof("[message: %s, reply: %s]", node.Children[0].NodeValue, reply)
				if len(reply) <= 0 {
					continue
				}

				i++
				//opts.addVideoCommentsForTopNReviewer(i+1, reply)
			}
		}

	}

	return nil
}

func (opts *CommentOption) addVideoCommentsForTopNReviewer(index int, reply string) error {
	var (
		video_clickReplySpanSelector = "#douyin-right-container > div:nth-child(2) > div > div.leftContainer.w0R6mo9z > div.tMWlo89q > div > div > div.sX7gMtFl.comment-mainContent > div:nth-child(%d) > div > div.RHiEl2d8 > div > div.rJFDwdFI > div > div.NRiH5zYV.ylYMpwSs"
		video_inputReplySpan         = "#douyin-right-container > div:nth-child(2) > div > div.leftContainer.w0R6mo9z > div.tMWlo89q > div > div > div.sX7gMtFl.comment-mainContent > div:nth-child(%d) > div > div.RHiEl2d8 > div > div.rJFDwdFI > div.yBfAEOf6 > div > div > div.pOcGUDMb > div > div > div.DraftEditor-editorContainer > div"
	)

	url1 := fmt.Sprintf(video_clickReplySpanSelector, index)
	url2 := fmt.Sprintf(video_inputReplySpan, index)

	newReply := []rune(reply)
	if len(newReply) > 100 {
		reply = string(newReply[0:95]) + "..."
	}
	if err := chromedp.Run(opts.ctx, chromedp.Click(url1), chromedp.Sleep(1*time.Second), chromedp.SendKeys(url2, reply+kb.Enter), chromedp.Sleep(5*time.Second)); err != nil {
		return err
	}
	return nil
}

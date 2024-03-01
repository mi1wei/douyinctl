package cmd

import (
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/mi1wei/douyinctl/config"
	"github.com/mi1wei/douyinctl/util/cookie"
	"github.com/mi1wei/douyinctl/util/file"
	"github.com/mi1wei/douyinctl/util/my_chromedp"
	"github.com/spf13/cobra"
)

type LoginOption struct {
	config *config.Config
}

var loginCmd = &cobra.Command{
	Use:     "login",
	Short:   "Login to platform",
	Example: "douyinctl login --platform douyin",
	Run: func(cmd *cobra.Command, args []string) {
		err := login()
		if err != nil {
			fmt.Println("Error: login failed, err: ", err)
			return
		}
		fmt.Println("Login successful")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

func login() error {
	opts := &LoginOption{
		config: &config.Config{},
	}

	if err := opts.config.Load(); err != nil {
		config.DefaultConfig.Store()
		return nil
	}
	return opts.performLogin()
}

func (opts *LoginOption) performLogin() error {
	loginUrl := "https://www.douyin.com"
	homePageUrl := "https://www.douyin.com/user/self"
	loginStatusSelector := "#island_b69f5 > div > div:nth-child(5) > div > a > div > img"
	defer opts.config.Store()

	for i, account := range opts.config.Accounts {
		ctx, cancel := my_chromedp.NewContext(false)
		defer cancel()

		if opts.config.Accounts[i].CookieFilePath != "" && file.IsExists("douyin", account.UserId) {
			if err := chromedp.Run(ctx, cookie.LoadCookies("douyin", account.UserId), chromedp.Navigate(loginUrl)); err != nil {
				return err
			}
		} else {
			if err := chromedp.Run(ctx, chromedp.Navigate(loginUrl)); err != nil {
				return err
			}
		}

		opts.config.Accounts[i].CookieFilePath = ""

		if err := chromedp.Run(ctx, chromedp.WaitVisible(loginStatusSelector)); err != nil {
			return err
		}

		if err := chromedp.Run(ctx, chromedp.Navigate(homePageUrl), chromedp.WaitVisible(loginStatusSelector), cookie.SaveCookie("douyin", account.UserId)); err != nil {
			return err
		}

		var name string
		var nameSelector = "#douyin-right-container > div.tQ0DXWWO.DAet3nqK.userNewUi > div > div > div.o1w0tvbC.F3jJ1P9_.InbPGkRv > div.mZmVWLzR > div.ds1zWG9C > h1 > span > span > span > span > span > span"
		if err := chromedp.Run(ctx, chromedp.Text(nameSelector, &name)); err != nil {
			return err
		}

		opts.config.Accounts[i].UserName = name
		opts.config.Accounts[i].CookieFilePath = fmt.Sprintf("./douyin/cookies_%s.json", account.UserId)
	}
	return nil
}

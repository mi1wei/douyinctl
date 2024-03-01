package cookie

import (
	"context"
	"errors"
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
)

func SaveCookie(platform string, userId string) chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		logging := logrus.WithField("platform", platform).WithField("userId", userId)
		cookieFilePath := fmt.Sprintf("./.%s/cookies_%s.json", platform, userId)
		cookies, err := network.GetCookies().Do(ctx)
		if err != nil {
			logging.Errorf("Failed to get cookies %v", err)
			return err
		}

		dir := filepath.Dir(cookieFilePath)
		if err = os.MkdirAll(dir, 0755); err != nil {
			logging.Errorf("Failed to create directory %v", err)
			return
		}

		cookiesData, err := network.GetCookiesReturns{Cookies: cookies}.MarshalJSON()
		if err != nil {
			logging.Errorf("Failed to serialize cookies %v", err)
			return err
		}

		if err = ioutil.WriteFile(cookieFilePath, cookiesData, 0755); err != nil {
			logging.Errorf("Failed to write cookies to file %v", err)
			return err
		}
		logging.Infof("Info: Cookies saved to %s", cookieFilePath)
		return
	}
}

var Enc = mahonia.NewEncoder("GBK")

func LoadCookies(platform string, userId string) chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		logging := logrus.WithField("platform", platform).WithField("userId", userId)
		cookieFilePath := fmt.Sprintf("./.%s/cookies_%s.json", platform, userId)
		if _, _err := os.Stat(cookieFilePath); os.IsNotExist(_err) {
			return errors.New("cookie file does not exist")
		}

		cookiesData, err := ioutil.ReadFile(cookieFilePath)
		if err != nil {
			logging.Errorf("Failed to read cookie file %v", err)
			return err
		}

		cookiesParams := network.SetCookiesParams{}
		if err = cookiesParams.UnmarshalJSON(cookiesData); err != nil {
			logging.Errorf("Failed to unmarshal cookie data %v", err)
			return err
		}

		err = network.SetCookies(cookiesParams.Cookies).Do(ctx)
		if err != nil {
			logging.WithField("platform", platform).Errorf("Failed to set cookies %v", err)
			return err
		}

		return nil
	}
}

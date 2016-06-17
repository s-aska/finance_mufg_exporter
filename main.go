package main

import (
	"github.com/labstack/gommon/log"
	"os"
	"sourcegraph.com/sourcegraph/go-selenium"
	"time"
)

const LOGIN_URL = "https://entry11.bk.mufg.jp/ibg/dfw/APLIN/loginib/login?_TRANID=AA000_001"

func main() {
	var webDriver selenium.WebDriver
	var err error
	caps := selenium.Capabilities(map[string]interface{}{"browserName": "chrome"})
	if webDriver, err = selenium.NewRemote(caps, os.Getenv("SELENIUM_URL")); err != nil {
		log.Errorf("Failed to open session: %s", err)
		return
	}
	defer webDriver.Quit()

	webDriver.SetImplicitWaitTimeout(5000)

	if err = webDriver.Get(LOGIN_URL); err != nil {
		log.Errorf("Failed to load page: url:%s error:%s", LOGIN_URL, err)
		return
	}

	if title, err := webDriver.Title(); err != nil {
		log.Errorf("Failed to get page title: %s", err)
		save(webDriver, "error")
		return
	} else {
		log.Infof("Page title: %s", title)
	}

	// ログイン
	if elem, err := webDriver.FindElement(selenium.ById, "account_id"); err != nil {
		log.Errorf("Failed to find element: id:%s error:%s", "account_id", err)
		save(webDriver, "error")
		return
	} else {
		elem.SendKeys(os.Getenv("MUFG_ID"))
	}

	if elem, err := webDriver.FindElement(selenium.ById, "ib_password"); err != nil {
		log.Errorf("Failed to find element: id:%s error:%s", "ib_password", err)
		save(webDriver, "error")
		return
	} else {
		elem.SendKeys(os.Getenv("MUFG_PASSWORD"))
	}

	if elem, err := webDriver.FindElement(selenium.ByClassName, "admb_m"); err != nil {
		log.Errorf("Failed to find element: className:%s error:%s", "admb_m", err)
		save(webDriver, "error")
		return
	} else {
		elem.Click()
		time.Sleep(2 * time.Second)
	}

	// ワンタイムパスワード
	if elem, err := webDriver.FindElement(selenium.ById, "ib_password"); err != nil {
		log.Infof("Don't need onetime login")
	} else {
		log.Infof("Need onetime login")
		elem.SendKeys(os.Getenv("MUFG_ONETIME"))
		if btn, err := webDriver.FindElement(selenium.ByCSSSelector, ".buttons a"); err != nil {
			log.Errorf("Failed to find element: selector:%s error:%s", ".buttons a", err)
			save(webDriver, "error")
			return
		} else {
			btn.Click()
			time.Sleep(2 * time.Second)
		}
	}

	// TODO: お知らせを既読にする

	// 残高取得
	if amount, err := webDriver.FindElement(selenium.ById, "setAmountDisplay"); err != nil {
		log.Errorf("Failed to find element: id:%s error:%s", "setAmountDisplay", err)
		save(webDriver, "error")
		return
	} else {
		if text, err := amount.Text(); err != nil {
			log.Error(err)
			save(webDriver, "error")
		} else {
			log.Infof(text)
		}
	}

	// ログアウト
	if btn, err := webDriver.FindElement(selenium.ByCSSSelector, ".logout a"); err != nil {
		log.Errorf("Failed to find element: selector:%s error:%s", ".buttons a", err)
		save(webDriver, "error")
		return
	} else {
		btn.Click()
		time.Sleep(2 * time.Second)
	}
}

func save(webDriver selenium.WebDriver, id string) {
	// Save Page title
	if title, err := webDriver.Title(); err == nil {
		log.Infof("Page title: %s", title)
	}

	// Save Page source
	if data, err := webDriver.PageSource(); err != nil {
		log.Error(err)
	} else {
		if fo, err := os.Create("output-" + id + ".html"); err != nil {
			log.Error(err)
		} else {
			fo.Write([]byte(data))
		}
	}

	// Save Screenshot
	if data, err := webDriver.Screenshot(); err != nil {
		log.Error(err)
	} else {
		if fo, err := os.Create("output-" + id + ".png"); err != nil {
			log.Error(err)
		} else {
			fo.Write(data)
		}
	}
}

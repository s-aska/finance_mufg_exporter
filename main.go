package main

import (
	"fmt"
	"os"
	"sourcegraph.com/sourcegraph/go-selenium"
)

func main() {
	var webDriver selenium.WebDriver
	var err error
	caps := selenium.Capabilities(map[string]interface{}{"browserName": "chrome"})
	if webDriver, err = selenium.NewRemote(caps, os.Getenv("SELENIUM_URL")); err != nil {
		fmt.Printf("Failed to open session: %s\n", err)
		return
	}
	defer webDriver.Quit()

	err = webDriver.Get("https://entry11.bk.mufg.jp/ibg/dfw/APLIN/loginib/login?_TRANID=AA000_001")
	if err != nil {
		fmt.Printf("Failed to load page: %s\n", err)
		return
	}

	if title, err := webDriver.Title(); err == nil {
		fmt.Printf("Page title: %s\n", title)
	} else {
		fmt.Printf("Failed to get page title: %s", err)
		return
	}

	setVal(webDriver, "account_id", os.Getenv("MUFG_ID"))
	setVal(webDriver, "ib_password", os.Getenv("MUFG_PASSWORD"))
	click(webDriver, "admb_m")

	// Onetime
	elem, err := webDriver.FindElement(selenium.ById, "onetime_password")
	if err == nil {
		println("Need Onetime")
		elem.SendKeys(os.Getenv("MUFG_ONETIME"))
		btn, err := webDriver.FindElement(selenium.ByCSSSelector, ".buttons a")
		if err != nil {
			fmt.Printf("Failed to find element: %s\n", err)
			return
		}
		btn.Click()
	}

	amount, err := webDriver.FindElement(selenium.ById, "setAmountDisplay")
	if err != nil {
		fmt.Printf("Failed to find element: %s\n", err)
		if title, err := webDriver.Title(); err == nil {
			fmt.Printf("Page title: %s\n", title)
		} else {
			fmt.Printf("Failed to get page title: %s", err)
			return
		}
		return
	}
	text, _ := amount.Text()
	println(text)
}

func setVal(webDriver selenium.WebDriver, id string, val string) {
	elem, err := webDriver.FindElement(selenium.ById, id)
	if err != nil {
		fmt.Printf("Failed to find element: %s\n", err)
		return
	}
	elem.SendKeys(val)
}

func click(webDriver selenium.WebDriver, id string) {
	elem, err := webDriver.FindElement(selenium.ByClassName, id)
	if err != nil {
		fmt.Printf("Failed to find element: %s\n", err)
		return
	}
	elem.Click()
}

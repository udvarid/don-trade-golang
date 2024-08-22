package controller

import (
	"html/template"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/udvarid/don-trade-golang/authenticator"
	"github.com/udvarid/don-trade-golang/collector"
	"github.com/udvarid/don-trade-golang/model"
)

var (
	activeConfiguration = &model.Configuration{}
)

func Init(config *model.Configuration) {
	activeConfiguration = config
	router := gin.Default()
	router.LoadHTMLGlob("html/*")
	// Serve static files from the "static" directory
	router.Static("/static", "./static")
	router.StaticFile("/favicon.ico", "./static/favicon.ico")

	router.GET("/", startPage)
	router.GET("/logout", logout)
	router.GET("/detailed/:item", detailedPage)
	router.POST("/validate/", validate)
	router.GET("/checkin/:id/:session", checkInTask)
	router.Run()
}

func logout(c *gin.Context) {
	id_cookie, err := c.Cookie("id")
	if err == nil {
		authenticator.Logout(id_cookie)
	}
	c.SetCookie("id", "", -1, "/", "localhost", false, true)
	c.SetCookie("session", "", -1, "/", "localhost", false, true)
	redirectTo(c, "/")
}

func validate(c *gin.Context) {
	var getSession GetSession
	c.BindJSON(&getSession)
	newSession, err := authenticator.GiveSession(getSession.Id)
	if err != nil {
		redirectTo(c, "/")
		return
	}

	isValidatedInTime := authenticator.Validate(activeConfiguration, getSession.Id, newSession, true)

	if isValidatedInTime {
		c.SetCookie("id", getSession.Id, 3600, "/", activeConfiguration.RemoteAddress, false, true)
		c.SetCookie("session", newSession, 3600, "/", activeConfiguration.RemoteAddress, false, true)
	}
	redirectTo(c, "/")
}

func detailedPage(c *gin.Context) {
	id := c.Param("item")
	items := collector.GetItemsFromItemMap(collector.GetItems())
	item := items[id]
	var pageCandle HtmlWithInfo
	html, _ := os.ReadFile("html/kline-detailed-" + id + ".html")
	pageCandle.Page = template.HTML(string(html))
	pageCandle.Name = item.Name
	pageCandle.Description = item.Description

	var pageCandle2 HtmlWithInfo
	html2, _ := os.ReadFile("html/kline-detailed2-" + id + ".html")
	pageCandle2.Page = template.HTML(string(html2))
	pageCandle2.Name = item.Name
	pageCandle2.Description = item.Description

	isLoggedIn := isLoggedIn(c)

	c.HTML(http.StatusOK, "detailed.html", gin.H{
		"title":         "Detailed Page",
		"isLoggedIn":    isLoggedIn,
		"description":   item.Description,
		"detailedPage1": pageCandle,
		"detailedPage2": pageCandle2,
	})
}

func startPage(c *gin.Context) {
	items := collector.GetItems()
	stockItems := items["stocks"]
	var stockPages []HtmlWithInfo
	for _, stockItem := range stockItems {
		html, _ := os.ReadFile("html/kline-" + stockItem.Name + ".html")
		var page HtmlWithInfo
		page.Page = template.HTML(string(html))
		page.Name = stockItem.Name
		page.Description = stockItem.Description
		stockPages = append(stockPages, page)

	}

	fxItems := items["fxs"]
	var fxPages []HtmlWithInfo
	for _, fxItem := range fxItems {
		html, _ := os.ReadFile("html/kline-" + fxItem.Name + ".html")
		var page HtmlWithInfo
		page.Page = template.HTML(string(html))
		page.Name = fxItem.Name
		page.Description = fxItem.Description
		fxPages = append(fxPages, page)

	}

	commodityItems := items["commodities"]
	var commodityPages []HtmlWithInfo
	for _, commodityItem := range commodityItems {
		html, _ := os.ReadFile("html/kline-" + commodityItem.Name + ".html")
		var page HtmlWithInfo
		page.Page = template.HTML(string(html))
		page.Name = commodityItem.Name
		page.Description = commodityItem.Description
		commodityPages = append(commodityPages, page)

	}

	cryptoItems := items["cryptos"]
	var cryptoPages []HtmlWithInfo
	for _, cryptoItem := range cryptoItems {
		html, _ := os.ReadFile("html/kline-" + cryptoItem.Name + ".html")
		var page HtmlWithInfo
		page.Page = template.HTML(string(html))
		page.Name = cryptoItem.Name
		page.Description = cryptoItem.Description
		cryptoPages = append(cryptoPages, page)

	}

	isLoggedIn := isLoggedIn(c)

	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":          "Main Page",
		"isLoggedIn":     isLoggedIn,
		"stockPages":     stockPages,
		"fxPages":        fxPages,
		"commodityPages": commodityPages,
		"cryptoPages":    cryptoPages,
	})

}

func isLoggedIn(c *gin.Context) bool {
	id_cookie, err := c.Cookie("id")
	isMissingCookie := false
	if err != nil {
		isMissingCookie = true
	}
	session_cookie, err := c.Cookie("session")
	if err != nil {
		isMissingCookie = true
	}
	if isMissingCookie {
		return false
	}
	return authenticator.IsValid(id_cookie, session_cookie)
}

func redirectTo(c *gin.Context, path string) {
	location := url.URL{Path: path}
	c.Redirect(http.StatusFound, location.RequestURI())
}

func checkInTask(c *gin.Context) {
	authenticator.CheckIn(c.Param("id"), c.Param("session"))
	c.HTML(http.StatusOK, "logged_in.html", gin.H{
		"title": "Logged in Page",
	})
}

type HtmlWithInfo struct {
	Name        string
	Description string
	Page        interface{}
}

type GetSession struct {
	Id string `json:"id"`
}

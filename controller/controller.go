package controller

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/udvarid/don-trade-golang/authenticator"
	chart "github.com/udvarid/don-trade-golang/chartBuilder"
	"github.com/udvarid/don-trade-golang/collector"
	"github.com/udvarid/don-trade-golang/model"
	"github.com/udvarid/don-trade-golang/orderManager"
	"github.com/udvarid/don-trade-golang/orderService"
	"github.com/udvarid/don-trade-golang/repository/candleRepository"
	userService "github.com/udvarid/don-trade-golang/user"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	activeConfiguration = &model.Configuration{}
)

func Init(config *model.Configuration) {
	activeConfiguration = config
	router := gin.Default()
	router.LoadHTMLGlob("html/*")

	router.Static("/static", "./static")
	router.StaticFile("/favicon.ico", "./static/favicon.ico")

	router.GET("/", startPage)
	router.GET("/logout", logout)
	router.GET("/user", user)
	router.POST("/addorder", addorder)
	router.GET("/deleteOrder/:order", deleteOrder)
	router.GET("/admin", admin)
	router.GET("/reset_db", resetDb)
	router.GET("/admin_order", adminOrder)
	router.GET("/detailed/:item", detailedPage)
	router.POST("/validate/", validate)
	router.GET("/checkin/:id/:session", checkInTask)
	router.Run()
}

func logout(c *gin.Context) {
	id_cookie, err := c.Cookie("id")
	_, session := getId(c)
	if err == nil {
		authenticator.Logout(id_cookie, session)
	}
	c.SetCookie("id", "", -1, "/", "localhost", false, true)
	c.SetCookie("session", "", -1, "/", "localhost", false, true)
	redirectTo(c, "/")
}

func deleteOrder(c *gin.Context) {
	isLoggedIn := isLoggedIn(c)
	if !isLoggedIn {
		redirectTo(c, "/")
	}

	if orderId, err := strconv.Atoi(c.Param("order")); err == nil {
		userId, _ := getId(c)
		orderService.DeleteOrder(orderId, userId)
	}

	redirectTo(c, "/user")
}

func user(c *gin.Context) {
	isLoggedIn := isLoggedIn(c)
	if !isLoggedIn {
		redirectTo(c, "/")
	}
	userId, session := getId(c)
	userStatistic := userService.GetUser(userId)
	candleSummary := candleRepository.GetAllCandleSummaries()[0]

	var pageBar HtmlWithInfo
	html, _ := os.ReadFile("html/kline-" + session + ".html")
	pageBar.Page = template.HTML(string(html))
	pageBar.Name = "Portfolio"
	pageBar.Description = "Detailed portfolio history for the user"

	orders := orderService.GetOrdersByUserId(userId)

	c.HTML(http.StatusOK, "user.html", gin.H{
		"title":         "user Page",
		"name":          userStatistic.Name,
		"assets":        transformUserAssetToString(userStatistic.Assets[:len(userStatistic.Assets)-1]),
		"totalAssets":   transformUserAssetToString(userStatistic.Assets[len(userStatistic.Assets)-1:])[0],
		"transactions":  transformTransactionToString(userStatistic.Transactions),
		"candleSummary": candleSummary.Summary,
		"barChart":      pageBar,
		"orders":        transformOrdersToString(orders),
	})
}

func admin(c *gin.Context) {
	isLoggedIn := isLoggedIn(c)
	userId, _ := getId(c)
	isAdminUser := isLoggedIn && activeConfiguration.Admin_user == userId
	if !isAdminUser {
		redirectTo(c, "/")
	}
	c.HTML(http.StatusOK, "admin.html", gin.H{
		"title": "Admin Page",
	})
}

func resetDb(c *gin.Context) {
	isLoggedIn := isLoggedIn(c)
	userId, _ := getId(c)
	isAdminUser := isLoggedIn && activeConfiguration.Admin_user == userId
	if !isAdminUser {
		redirectTo(c, "/")
	}
	collector.DeletePriceDatabase(activeConfiguration)
}

func adminOrder(c *gin.Context) {
	isLoggedIn := isLoggedIn(c)
	userId, _ := getId(c)
	isAdminUser := isLoggedIn && activeConfiguration.Admin_user == userId
	if !isAdminUser {
		redirectTo(c, "/")
	}
	orderManager.ServeOrders(false, userId)
}

func addorder(c *gin.Context) {
	isLoggedIn := isLoggedIn(c)
	if !isLoggedIn {
		redirectTo(c, "/")
	}
	userId, _ := getId(c)
	var order model.OrderInString
	c.BindJSON(&order)
	orderService.ValidateAndAddOrder(order, userId)
	redirectTo(c, "/user")
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
		chart.BuildUserHistoryChart(userService.GetUserHistory(getSession.Id, 30), newSession)
		redirectTo(c, "/")
	}
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
		"name":          item.Name,
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
	userId, _ := getId(c)
	isAdminUser := isLoggedIn && activeConfiguration.Admin_user == userId

	priceChanges := userService.GetPriceChanges()

	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":          "Main Page",
		"isLoggedIn":     isLoggedIn,
		"isAdminUser":    isAdminUser,
		"stockPages":     stockPages,
		"fxPages":        fxPages,
		"commodityPages": commodityPages,
		"cryptoPages":    cryptoPages,
		"priceChanges":   priceChanges,
	})

}

func checkInTask(c *gin.Context) {
	authenticator.CheckIn(c.Param("id"), c.Param("session"))
	c.HTML(http.StatusOK, "logged_in.html", gin.H{
		"title": "Logged in Page",
	})
}

func getId(c *gin.Context) (string, string) {
	id_cookie, _ := c.Cookie("id")
	session_cookie, _ := c.Cookie("session")
	return id_cookie, session_cookie
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

func transformTransactionToString(transactions []model.Transaction) []model.TransactionWitString {
	var result []model.TransactionWitString
	p := message.NewPrinter(language.Hungarian)
	for _, tr := range transactions {
		var trWithStr model.TransactionWitString
		trWithStr.Asset = tr.Asset
		trWithStr.Date = tr.Date.Format("06-01-02")
		trWithStr.Volume = p.Sprintf("%d", int(tr.Volume))
		if tr.Volume > 0 {
			trWithStr.Volume = "+" + trWithStr.Volume
		}
		result = append(result, trWithStr)
	}
	return result
}

func transformOrdersToString(orders []model.Order) []model.OrderInString {
	var result []model.OrderInString
	for _, order := range orders {
		var orderInString model.OrderInString
		orderInString.ID = order.ID
		orderInString.UserID = order.UserID
		orderInString.Item = order.Item
		orderInString.Direction = order.Direction
		orderInString.Type = order.Type
		if math.Abs(order.LimitPrice) < 0.0001 {
			orderInString.LimitPrice = "-"
		} else {
			orderInString.LimitPrice = fmt.Sprintf("%.1f", order.LimitPrice)
		}
		if math.Abs(order.NumberOfItems) < 0.0001 {
			orderInString.NumberOfItems = "-"
		} else {
			orderInString.NumberOfItems = fmt.Sprintf("%.1f", order.NumberOfItems)
		}
		if math.Abs(order.Usd) < 0.0001 {
			orderInString.Usd = "-"
		} else {
			orderInString.Usd = fmt.Sprintf("%.1f", order.Usd)
		}
		orderInString.AllIn = order.AllIn
		orderInString.ValidDays = strconv.Itoa(order.ValidDays)
		result = append(result, orderInString)
	}

	return result
}

func transformUserAssetToString(assets []model.AssetWithValue) []model.AssetWithValueInString {
	var result []model.AssetWithValueInString
	p := message.NewPrinter(language.Hungarian)
	for _, asset := range assets {
		var newAsset model.AssetWithValueInString
		newAsset.Item = asset.Item
		newAsset.Volume = p.Sprintf("%d", int(asset.Volume))
		newAsset.Price = fmt.Sprintf("%.2f", asset.Price)
		newAsset.Value = p.Sprintf("%d", int(asset.Value))
		result = append(result, newAsset)
	}
	return result
}

type HtmlWithInfo struct {
	Name        string
	Description string
	Page        interface{}
}

type GetSession struct {
	Id string `json:"id"`
}

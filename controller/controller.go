package controller

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/udvarid/don-trade-golang/authenticator"
	chart "github.com/udvarid/don-trade-golang/chartBuilder"
	"github.com/udvarid/don-trade-golang/collector"
	"github.com/udvarid/don-trade-golang/model"
	"github.com/udvarid/don-trade-golang/orderManager"
	"github.com/udvarid/don-trade-golang/orderService"
	"github.com/udvarid/don-trade-golang/repository/candleRepository"
	userService "github.com/udvarid/don-trade-golang/user"
	userstatistic "github.com/udvarid/don-trade-golang/userStatistic"
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
	router.GET("/transactions", transactions)
	router.GET("/users", users)
	router.GET("/user_settings", userSettings)
	router.GET("/user_delete", userDelete)
	router.GET("/clear_item/:item/:short", clearItem)
	router.POST("/addorder", addorder)
	router.POST("/modify_order", modifyOrder)
	router.GET("/deleteOrder/:order", deleteOrder)
	router.GET("/admin", admin)
	router.GET("/reset_db", resetDb)
	router.GET("/admin_order", adminOrder)
	router.GET("/detailed/:item", detailedPage)
	router.POST("/validate/", validate)
	router.POST("/name_change/", nameChange)
	router.POST("/notify_change/", notifyChange)
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

func clearItem(c *gin.Context) {
	isLoggedIn := isLoggedIn(c)
	if !isLoggedIn {
		redirectTo(c, "/")
	}
	userId, _ := getId(c)
	shortParam := c.Param("short")
	shortBool, err := strconv.ParseBool(shortParam)
	if err == nil {
		orderService.MakeClearOrder(userId, c.Param("item"), shortBool)
	}
	redirectTo(c, "/user")
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

func userDelete(c *gin.Context) {
	isLoggedIn := isLoggedIn(c)
	if !isLoggedIn {
		redirectTo(c, "/")
	}
	id, session := getId(c)
	authenticator.Logout(id, session)
	userService.DeleteUser(id)
	c.SetCookie("id", "", -1, "/", "localhost", false, true)
	c.SetCookie("session", "", -1, "/", "localhost", false, true)
	redirectTo(c, "/")
}

func userSettings(c *gin.Context) {
	isLoggedIn := isLoggedIn(c)
	if !isLoggedIn {
		redirectTo(c, "/")
	}
	userId, _ := getId(c)
	user := userService.GetUser(userId)
	c.HTML(http.StatusOK, "user_settings.html", gin.H{
		"title":               "user_settings Page",
		"name":                user.Name,
		"notifyAtTransaction": user.Config.NotifyAtTransaction,
		"notifyDaily":         user.Config.NotifyDaily,
	})
}

func users(c *gin.Context) {
	traders := userService.GetTraders()
	c.HTML(http.StatusOK, "users.html", gin.H{
		"title":   "users Page",
		"traders": transformTradersToString(traders),
	})
}

func transactions(c *gin.Context) {
	isLoggedIn := isLoggedIn(c)
	if !isLoggedIn {
		redirectTo(c, "/")
	}
	userId, _ := getId(c)
	userStatistic := userstatistic.GetUserStatistic(userId, true)

	c.HTML(http.StatusOK, "transactions.html", gin.H{
		"title":        "Transactions Page",
		"name":         userStatistic.Name,
		"transactions": transformTransactionToString(userStatistic.Transactions),
	})
}

func user(c *gin.Context) {
	isLoggedIn := isLoggedIn(c)
	if !isLoggedIn {
		redirectTo(c, "/")
	}
	userId, session := getId(c)
	userStatistic := userstatistic.GetUserStatistic(userId, false)
	candleSummary := candleRepository.GetAllCandleSummaries()[0]

	var charts []HtmlWithInfo
	var pageBar HtmlWithInfo
	html, _ := os.ReadFile("html/kline-" + session + ".html")
	pageBar.Page = template.HTML(string(html))
	pageBar.Name = "Total"
	pageBar.Description = "Detailed portfolio history for the user"
	charts = append(charts, pageBar)
	for _, asset := range userStatistic.Assets[:len(userStatistic.Assets)-2] {
		html_asset, _ := os.ReadFile("html/kline-" + asset.Item + ".html")
		htmlContent := string(html_asset)
		htmlContent = strings.Replace(htmlContent, "width:900px", "width:750px", 1)
		var assetChart HtmlWithInfo
		assetChart.Page = template.HTML(htmlContent)
		assetChart.Name = asset.Item
		charts = append(charts, assetChart)
	}

	orders := orderService.GetOrdersByUserId(userId)

	c.HTML(http.StatusOK, "user.html", gin.H{
		"title":         "user Page",
		"name":          userStatistic.Name,
		"assets":        transformUserAssetToString(userStatistic.Assets[:len(userStatistic.Assets)-2]),
		"usd":           transformUserAssetToString(userStatistic.Assets[len(userStatistic.Assets)-2 : len(userStatistic.Assets)-1])[0],
		"totalAssets":   transformUserAssetToString(userStatistic.Assets[len(userStatistic.Assets)-1:])[0],
		"candleSummary": candleSummary.Summary,
		"creditLimit":   fmt.Sprintf("%.1f", userStatistic.CreditLimit*100) + "%",
		"charts":        charts,
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

func modifyOrder(c *gin.Context) {
	isLoggedIn := isLoggedIn(c)
	if !isLoggedIn {
		redirectTo(c, "/")
	}
	userId, _ := getId(c)
	var order model.OrderModifyInString
	c.BindJSON(&order)
	orderService.ModifyOrder(userId, &order)
	redirectTo(c, "/user")
}

func addorder(c *gin.Context) {
	isLoggedIn := isLoggedIn(c)
	if !isLoggedIn {
		redirectTo(c, "/")
	}
	userId, _ := getId(c)
	var order model.OrderInString
	c.BindJSON(&order)
	orderService.ValidateAndAddOrder(&order, userId)
	redirectTo(c, "/user")
}

func notifyChange(c *gin.Context) {
	isLoggedIn := isLoggedIn(c)
	if !isLoggedIn {
		redirectTo(c, "/")
	}
	userId, _ := getId(c)
	var newNotify NotifyChange
	c.BindJSON(&newNotify)
	userService.ChangeNotify(userId, newNotify.Transaction, newNotify.Daily)
	redirectTo(c, "/user")
}

func nameChange(c *gin.Context) {
	isLoggedIn := isLoggedIn(c)
	if !isLoggedIn {
		redirectTo(c, "/")
	}
	userId, _ := getId(c)
	var newName GetSession
	c.BindJSON(&newName)
	if newName.Id != "" {
		userService.ChangeName(userId, newName.Id)
	}
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
		sessionTime := 3600 * 24 * 4
		c.SetCookie("id", getSession.Id, sessionTime, "/", activeConfiguration.RemoteAddress, false, true)
		c.SetCookie("session", newSession, sessionTime, "/", activeConfiguration.RemoteAddress, false, true)
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
		"items":         items,
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

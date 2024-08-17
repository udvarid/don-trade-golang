package controller

import (
	"html/template"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/udvarid/don-trade-golang/collector"
)

func Init() {
	router := gin.Default()
	router.LoadHTMLGlob("html/*")

	router.GET("/", startPage)
	router.GET("/detailed/:item", detailedPage)
	router.Run()
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
	c.HTML(http.StatusOK, "detailed.html", gin.H{
		"title":         "Detailed Page",
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

	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":          "Main Page",
		"stockPages":     stockPages,
		"fxPages":        fxPages,
		"commodityPages": commodityPages,
		"cryptoPages":    cryptoPages,
	})

}

type HtmlWithInfo struct {
	Name        string
	Description string
	Page        interface{}
}

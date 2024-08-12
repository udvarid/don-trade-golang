package controller

import (
	"html/template"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func Init() {
	router := gin.Default()
	router.LoadHTMLGlob("html/*")

	router.GET("/", startPage)
	router.Run()
}

func startPage(c *gin.Context) {
	// Read the pre-generated HTML files
	html1, _ := os.ReadFile("html/kline-AMZN.html")
	html2, _ := os.ReadFile("html/kline-BTCUSD.html")
	html3, _ := os.ReadFile("html/kline-EURUSD.html")
	html4, _ := os.ReadFile("html/kline-CLUSD.html")

	data := map[string]interface{}{
		"Html1": template.HTML(string(html1)),
		"Html2": template.HTML(string(html2)),
		"Html3": template.HTML(string(html3)),
		"Html4": template.HTML(string(html4)),
	}

	/*c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Home Page",
		"pages": data,
	})*/

	c.HTML(http.StatusOK, "index.html", data)

}

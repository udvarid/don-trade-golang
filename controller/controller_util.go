package controller

import (
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/udvarid/don-trade-golang/authenticator"
	"github.com/udvarid/don-trade-golang/collector"
	"github.com/udvarid/don-trade-golang/model"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

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
			orderInString.LimitPrice = fmt.Sprintf("%.2f", order.LimitPrice)
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
		orderInString.Short = order.Short
		orderInString.ValidDays = strconv.Itoa(order.ValidDays)
		result = append(result, orderInString)
	}

	return result
}

func transformTradersToString(traders []model.UserSummary) []model.UserSummaryInString {
	var result []model.UserSummaryInString
	for _, trader := range traders {
		var traderInString model.UserSummaryInString
		traderInString.UserID = trader.UserID
		traderInString.UserName = trader.UserName
		traderInString.Profit = fmt.Sprintf("%.2f", trader.Profit*100) + "%"
		traderInString.TraderSince = trader.TraderSince
		traderInString.Invested = fmt.Sprintf("%.2f", trader.Invested*100) + "%"
		traderInString.CreditLimit = fmt.Sprintf("%.1f", trader.CreditLimit*100) + "%"
		result = append(result, traderInString)
	}
	return result
}

func transformUserAssetToString(assets []model.AssetWithValue) []model.AssetWithValueInString {
	var result []model.AssetWithValueInString
	p := message.NewPrinter(language.Hungarian)
	items := collector.GetItemsFromItemMap(collector.GetItems())
	for _, asset := range assets {
		var newAsset model.AssetWithValueInString
		if asset.Item == "USD" || asset.Item == "Total" {
			newAsset.Item = asset.Item
		} else {
			newAsset.Item = asset.Item + " - " + items[asset.Item].Description
		}
		newAsset.ItemPure = asset.Item
		newAsset.Volume = p.Sprintf("%d", int(asset.Volume))
		newAsset.Price = fmt.Sprintf("%.2f", asset.Price)
		newAsset.Value = p.Sprintf("%d", int(asset.Value))
		if math.Abs(asset.BookValue) > 0.0001 {
			newAsset.BookValue = p.Sprintf("%d", int(asset.BookValue))
			if asset.BookValue > 0.0001 {
				newAsset.Profit = fmt.Sprintf("%.2f", (asset.Value/asset.BookValue-1)*100) + "%"
				newAsset.Short = false
			} else {
				newAsset.Profit = fmt.Sprintf("%.2f", (asset.Value/asset.BookValue-1)*-100) + "%"
				newAsset.Short = true
			}

		} else {
			newAsset.BookValue = "-"
			newAsset.Profit = "-"
		}
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

type NotifyChange struct {
	Transaction bool `json:"transaction"`
	Daily       bool `json:"daily"`
}

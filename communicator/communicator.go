package communicator

import (
	"log"
	"net/smtp"

	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/udvarid/don-trade-golang/model"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	activeConfiguration = &model.Configuration{}
)

func Init(config *model.Configuration) {
	activeConfiguration = config
}

var verifier = emailverifier.NewVerifier()

func SendMessageAboutStatus(userStatistic *model.UserStatistic) {
	ret, err := verifier.Verify(userStatistic.ID)
	if err == nil && ret.Syntax.Valid {
		p := message.NewPrinter(language.Hungarian)
		portfolioString := ""
		for _, asset := range userStatistic.Assets {
			if asset.Item != "Total" {
				portfolioString += "<p>" + asset.Item + " - Volume:" + p.Sprintf("%d", int(asset.Volume)) + " Value: " + p.Sprintf("%d", int(asset.Value)) + "</p>"
			} else {
				portfolioString += "<p> ---------------------------------------</p>"
				portfolioString += "<p>" + asset.Item + " - Value: " + p.Sprintf("%d", int(asset.Value)) + "</p>"
			}
		}
		msg := []byte("To: " + userStatistic.ID + "\r\n" +
			"Subject: Daily Status from Don-Trade!\r\n" +
			"MIME-version: 1.0;\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
			"\r\n" +
			"<html><body>" +
			"<p>Wellcome " + userStatistic.Name + "!</p>" +
			"<p></p>" +
			"<p>Here is your short summary of your portfolio:</p>" +
			portfolioString +
			"<p></p>" +
			"<p><a href=" + activeConfiguration.RemoteAddress + ">Visit us</a> soon to make more profit!:) </p>" +
			"</body></html>")
		sendMail(userStatistic.ID, msg)
	}
}

func SendMessageAboutOrders(toAddress string, orders []model.CompletedOrderToMail) {
	ret, err := verifier.Verify(toAddress)
	if err == nil && ret.Syntax.Valid {
		orderString := ""
		for _, order := range orders {
			orderString += "<p>" + order.Item + " - " + order.Type + ". Volumen:" + order.Volumen + " Price:" + order.Price + " Usd: " + order.Usd + "</p>"
		}
		msg := []byte("To: " + toAddress + "\r\n" +
			"Subject: Completed orders!\r\n" +
			"MIME-version: 1.0;\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
			"\r\n" +
			"<html><body>" +
			"<p>Your following orders have been executed:</p>" +
			orderString +
			"</body></html>")
		sendMail(toAddress, msg)
	}
}

func SendMessageWithLink(toAddress string, toLink string) {
	ret, err := verifier.Verify(toAddress)
	if err == nil && ret.Syntax.Valid {
		msg := []byte("To: " + toAddress + "\r\n" +
			"Subject: Please check in!\r\n" +
			"MIME-version: 1.0;\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
			"\r\n" +
			"<html><body>" +
			"<p>To login please: <a href=" + toLink + ">click me</a></p>" +
			"</body></html>")
		sendMail(toAddress, msg)
	}
}

func sendMail(toAddress string, message []byte) {
	address, pw := getMailAndPw()
	auth := smtp.PlainAuth("", address, pw, "smtp.gmail.com")
	err := smtp.SendMail("smtp.gmail.com:587", auth, address, []string{toAddress}, message)
	if err != nil {
		log.Print(err)
	}
}

func getMailAndPw() (string, string) {
	var address string
	var pw string
	if activeConfiguration.Environment == "local" {
		address = activeConfiguration.Mail_local_from
		pw = activeConfiguration.Mail_local_psw
	} else {
		address = activeConfiguration.Mail_from
		pw = activeConfiguration.Mail_psw
	}
	return address, pw
}

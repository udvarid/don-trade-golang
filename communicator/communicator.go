package communicator

import (
	"log"
	"net/smtp"

	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/udvarid/don-trade-golang/model"
)

var (
	activeConfiguration = &model.Configuration{}
)

func Init(config *model.Configuration) {
	activeConfiguration = config
}

var verifier = emailverifier.NewVerifier()

func SendMessageAboutOrders(toAddress string, orders []model.CompletedOrderToMail) {
	ret, err := verifier.Verify(toAddress)
	orderString := ""
	for _, order := range orders {
		orderString += "<p>" + order.Item + " - " + order.Type + ". Volumen:" + order.Volumen + " Price:" + order.Price + " Usd: " + order.Usd + "</p>"
	}
	if err == nil && ret.Syntax.Valid {
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

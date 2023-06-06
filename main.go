package main

import (
	"log"
	"fmt"
	"time"
	"os"
	"strings"
	"image/png"
	"github.com/mymmrac/telego"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/joho/godotenv"
	tu "github.com/mymmrac/telego/telegoutil"
	"bytes"
	"./squareapi"
)


func sendQrCode(bot *telego.Bot, invoiceId string, orderId string, chatId telego.ChatID) {
	paymentLink := os.Getenv("SQUARE_PAYMENT_URL") + invoiceId
	log.Printf("Invoice: %s", invoiceId)
	log.Printf("Invoice: %s", orderId)
	log.Printf("Payment link: %s", paymentLink)
	log.Printf("ChatID: %s", chatId)
	qrCode, _ := qr.Encode(paymentLink, qr.M, qr.Auto)
	qrCode, _ = barcode.Scale(qrCode, 250, 250)
	image := new(bytes.Buffer)
	png.Encode(image, qrCode)
	photo := tu.Photo(
		chatId,
		tu.File(tu.NameReader(image, "Invoice")),
	).WithCaption(fmt.Sprintf("Invoice for order \"%s\".\nLink: %s", orderId, paymentLink))
	bot.SendPhoto(photo)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error occured while reading env file: %s.", err)
	}
	bot, err := telego.NewBot(os.Getenv("BOT_TOKEN"), telego.WithDefaultDebugLogger())
	if err != nil {
		log.Fatalf("Err: %s", err)
		os.Exit(1)
	}

	botUser, err := bot.GetMe()

	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	log.Printf("Bot user: $+v", botUser)

	updates, _ := bot.UpdatesViaLongPolling(nil)

	defer bot.StopLongPolling()

	var api = squareapi.SquareApi{
		LocationId: os.Getenv("SELLER_LOCATION_ID"),
		Version: os.Getenv("SQUARE_API_VERSION"),
		Endpoint: os.Getenv("SQUARE_API_ENDPOINT"),
		AccessToken: os.Getenv("SQUARE_ACCESS_TOKEN"),
	}

	for update := range updates {
		if update.Message != nil {
			chatId := tu.ID(update.Message.Chat.ID)

			bot.SendMessage(tu.Message(chatId, "Let me think..."))
			// _, _ = bot.CopyMessage(tu.CopyMessage(chatId, chatId, update.Message.MessageID),)
			commandSplit := strings.Split(update.Message.Text, " ")
			switch command := commandSplit[0]; command {
			case "/invoice":
			  // var orderId = "Ul3Or1Q1QrwaF6XgzhgkeIsNua4F"
				var orderId = commandSplit[1]
				var customerId = commandSplit[2]
				// var invoiceTitle = commandSplit[3]
			  var invoiceTitle = "Customer invoice"
			  var defaultInvoiceDueDate = time.Now().AddDate(0, 1, 0)
			  dueDate := fmt.Sprintf("%d-%02d-%02d",
				  defaultInvoiceDueDate.Year(),
				  defaultInvoiceDueDate.Month(),
				  defaultInvoiceDueDate.Day())
			  invoice, err := api.CreateInvoice(orderId, customerId, invoiceTitle, dueDate)
			  if err != "" {
				  invoices, _ := api.SearchInvoices()
				  for _, invoice := range invoices {
					  if invoice.OrderId == orderId {
						  sendQrCode(bot, invoice.Id, invoice.OrderId, chatId)

					  }
				  }
				  return

			  }
			  sendQrCode(bot, invoice.Id, invoice.OrderId, chatId)
			case "/orders":
				orders, _ := api.GetOrders()
				for _, order := range orders {
					bot.SendMessage(tu.Message(chatId, order.Id))
				}
			case "/customers":
				customers, _ := api.GetCustomers()
				for _, customer := range customers {
					bot.SendMessage(tu.Message(chatId, customer.Id))
				}
			default:
				bot.SendMessage(tu.Message(chatId, "Unknown command."))
			}

		}

		if update.CallbackQuery != nil {
			log.Print("Received callback with data: ", update.CallbackQuery.Data)
		}
	}
}

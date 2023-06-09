package main

import (
	"bot/squareapi"
	"bytes"
	"fmt"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"image/png"
	"log"
	"os"
	"strings"
	"time"
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
		LocationId:  os.Getenv("SELLER_LOCATION_ID"),
		Version:     os.Getenv("SQUARE_API_VERSION"),
		Endpoint:    os.Getenv("SQUARE_API_ENDPOINT"),
		AccessToken: os.Getenv("SQUARE_ACCESS_TOKEN"),
	}

	for update := range updates {
		if update.Message != nil {
			chatId := tu.ID(update.Message.Chat.ID)

			if !strings.HasPrefix(update.Message.Text, "/start") {
				bot.SendMessage(tu.Message(chatId, "Let me think..."))
			}
			// _, _ = bot.CopyMessage(tu.CopyMessage(chatId, chatId, update.Message.MessageID),)
			commandSplit := strings.Split(update.Message.Text, " ")
			switch command := commandSplit[0]; command {
			case "/invoice":
				// var orderId = "Ul3Or1Q1QrwaF6XgzhgkeIsNua4F"
				if len(commandSplit) != 3 {
					bot.SendMessage(tu.Message(chatId, "Command should be of the form: /invoice <order_id> <customer_id>"))
					continue
				}
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

				} else {
					sendQrCode(bot, invoice.Id, invoice.OrderId, chatId)
				}
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
			case "/start":
				bot.SendMessage(tu.Message(chatId, "Welcome! See the menu to get a list of commands."))
			default:
				bot.SendMessage(tu.Message(chatId, "Unknown command."))
			}

		}

		if update.CallbackQuery != nil {
			log.Print("Received callback with data: ", update.CallbackQuery.Data)
		}
	}
}

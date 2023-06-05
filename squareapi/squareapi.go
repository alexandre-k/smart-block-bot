package squareapi

import (
	"log"
	"io"
	"github.com/google/uuid"
	"net/http"
	"bytes"
	"encoding/json"
)


type PrimaryRecipientReq struct {
	CustomerId string `json:"customer_id"`
}


type PrimaryRecipient struct {
	CustomerId string `json:"customer_id"`
	EmailAddress string `json:"email_address"`
	FamilyName string `json:"family_name"`
	GivenName string `json:"given_name"`
	PhoneNumber string `json:"phone_number"`
}


type AmountMoney struct {
	Amount int `json:"amount"`
	Currency string `json:"currency"`
}

type PaymentRequestReq struct {
	RequestType string `json:"request_type"`
	DueDate string `json:"due_date"`
	AutomaticPaymentSource string `json:"automatic_payment_source"`
}

type PaymentRequest struct {
	Uid string `json:"uid"`
	RequestType string `json:"request_type"`
	DueDate string `json:"due_date"`
	TippingEnabled bool `json:"tipping_enabled"`
	ComputedAmountMoney AmountMoney `json:"amount_money"`
	TotalCompletedAmountMoney AmountMoney `json:"total_completed_amount_money"`
	AutomaticPaymentSource string `json:"automatic_payment_source"`
}


type AcceptedPaymentMethods struct {
	BankAccount bool `json:"bank_account"`
	BuyNowPayLater bool `json:"buy_now_pay_later"`
	Card bool `json:"card"`
	SquareGiftCard bool `json:"square_gift_card"`
}


type InvoiceReq struct {
	LocationId string `json:"location_id"`
	OrderId string `json:"order_id"`
	PrimaryRecipient PrimaryRecipientReq `json:"primary_recipient"`
	DeliveryMethod string `json:"delivery_method"`
	PaymentRequests []PaymentRequestReq `json:"payment_requests"`
	AcceptedPaymentMethods AcceptedPaymentMethods `json:"accepted_payment_methods"`
	Title string `json:"title"`
}


type Invoice struct {
	Id string `json:"id"`
	Version int `json:"version"`
	LocationId string `json:"location_id"`
	OrderId string `json:"order_id"`
	PaymentRequests []PaymentRequest `json:"payment_requests"`
	PrimaryRecipient PrimaryRecipient `json:"primary_recipient"`
	InvoiceNumber string `json:"invoice_number"`
	Title string `json:"title"`
	Status string `json:"status"`
	Timezone bool `json:"timezone"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	AcceptedPaymentMethods AcceptedPaymentMethods `json:"accepted_payment_methods"`
	DeliveryMethod string `json:"delivery_method"`
	StorePaymentMethodEnabled bool `json:"store_payment_method_enabled"`
	SaleOrServiceDate string `json:"sale_or_service_date"`
}


type ErrorResponse struct {
	Errors []Error `json:"errors"`
}
type Error struct {
	Category string `json:"category"`
	Code string `json:"code"`
	Detail string `json:"detail"`
}


type InvoiceResponse struct {
	Invoice Invoice `json:"invoice"`
}


type InvoicesResponse struct {
	Invoices []Invoice `json:"invoices"`
}


type LineItem struct {
	Uid string `json:"uid"`
	CatalogObjectId string `json:"catelog_object_id"`
	CatalogVersion string `json:"catalog_version"`
	Quantity string `json:"quantity"`
	name string `json:"name"`
	VariationName string `json:"variation_name"`
	BasePriceMoney AmountMoney `json:"base_price_money"`
	GrossSalesMoney AmountMoney `json:"gross_sales_money"`
	TotalTaxMoney AmountMoney `json:"total_tax_money"`
	TotalDiscountMoney AmountMoney `json:"total_discount_money"`
	TotalMoney AmountMoney `json:"total_money"`
	VariationTotalPriceMoney AmountMoney `json:"variation_total_price_money"`
	ItemType string `json:"item_type"`
	TotalServiceChargeMoney AmountMoney `json:"total_service_charge_money"`
}

type NetAmounts struct {
	TotalMoney AmountMoney `json:"total_money"`
	TaxMoney AmountMoney `json:"tax_money"`
	DiscountMoney AmountMoney `json:"discount_money"`
	TipMoney AmountMoney `json:"tip_money"`
	ServiceChargeMoney AmountMoney `json:"service_charge_money"`
}

type Source struct {
	Name string `json:"name"`
}


type Order struct {
	Id string `json:"id"`
	LocationId string `json:"location_id"`
	LineItems []LineItem `json:"line_items"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	State string `json:"state"`
	Version int `json:"version"`
	TotalTaxMoney AmountMoney `json:"total_tax_money"`
	TotalDiscountMoney AmountMoney `json:"total_discount_money"`
	TotalTipMoney AmountMoney `json:"total_tip_money"`
	TotalMoney AmountMoney `json:"total_money"`
	TotalServiceChargeMoney AmountMoney `json:"total_service_charge_money"`
	NetAmounts NetAmounts `json:"net_amounts`
	Source Source `json:"source"`
	CustomerId string `json:"customer_id"`
	NetAmountDueMoney AmountMoney `json:"net_amount_due_money"`
}


type OrdersResponse struct {
	Orders []Order `json:"orders"`
}


type Preferences struct {
	EmailUnsubscribed bool `json:"email_unsubscribed"`
}


type Customer struct {
	Id string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	GivenName string `json:"given_name"`
	FamilyName string `json:"family_name"`
	EmailAddress string `json:"email_address"`
	PhoneNumber string `json:"phone_number"`
	Preferences Preferences `json:"preferences"`
	CreationSource string `json:"creation_source"`
	Version int `json:"version"`
}


type CustomersResponse struct {
	Customers []Customer `json:"customers"`
}


type Params struct {
	Method string
	Endpoint string
	Data []byte
}


type SquareApi struct {
	LocationId string
	Version string
	AccessToken string
	Endpoint string

}


func (api SquareApi) createSquareRequest(params Params) ([]byte, int) {
	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, params.Data, "", "\t")
	log.Println(prettyJSON.String())
	client := &http.Client{}
	req, err := http.NewRequest(params.Method, api.Endpoint + params.Endpoint, bytes.NewBuffer(params.Data))
	req.Header.Set("Square-Version", api.Version)
	req.Header.Set("Authorization", "Bearer " + api.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	return body, res.StatusCode
}

func (api SquareApi) GetOrders() ([]Order, string) {
	queryFilter := map[string][]string{
		"location_ids": []string{api.LocationId}}
	bQueryFilter, _ := json.Marshal(queryFilter)
	body, statusCode := api.createSquareRequest(Params{Method: http.MethodPost, Endpoint: "/orders/search", Data: bQueryFilter})
	// decoder := json.NewDecoder(bytes.NewReader(body))
	if statusCode >= 400 {
		// var mapping map[string]interface{}
		// decoder.Decode(&mapping)
		var data ErrorResponse
		json.Unmarshal(body, &data)
		log.Println("Error: ", data.Errors)
		return []Order{}, data.Errors[0].Detail
	} else {
		// data := InvoiceResponse{}
		// var mapping map[string]interface{}
		var mapping OrdersResponse
		// err := decoder.Decode(&mapping)
		json.Unmarshal(body, &mapping)
		log.Println(mapping)

		return mapping.Orders, ""
	}
}


func (api SquareApi) GetCustomers() ([]Customer, string) {
	body, statusCode := api.createSquareRequest(Params{Method: http.MethodGet, Endpoint: "/customers"})
	if statusCode >= 400 {
		var data ErrorResponse
		json.Unmarshal(body, &data)
		log.Println("Error: ", data.Errors)
		return []Customer{}, data.Errors[0].Detail
	} else {
		var mapping CustomersResponse
		json.Unmarshal(body, &mapping)
		log.Println(mapping)

		return mapping.Customers, ""
	}
}


func (api SquareApi) SearchInvoices() ([]Invoice, string) {
	queryFilter := map[string]map[string]map[string][]string{
		"query": {"filter": {"location_ids": []string{api.LocationId}}}}
	bQueryFilter, _ := json.Marshal(queryFilter)
	body, statusCode := api.createSquareRequest(Params{Method: http.MethodPost, Endpoint: "/invoices/search", Data: bQueryFilter})
	if statusCode >= 400 {
		var data ErrorResponse
		json.Unmarshal(body, &data)
		log.Println("Error: ", data.Errors)
		return []Invoice{}, data.Errors[0].Detail
	} else {
		var mapping InvoicesResponse
		json.Unmarshal(body, &mapping)
		log.Println(mapping)

		return mapping.Invoices, ""
	}
}

func (api SquareApi) CreateInvoice(
	orderId string,
	customerId string,
	invoiceTitle string,
	dueDate string) (Invoice, string) {
	invoice := InvoiceReq{
		LocationId: api.LocationId,
		OrderId: orderId,
		PrimaryRecipient: PrimaryRecipientReq{
			CustomerId: customerId,
		},
		PaymentRequests: []PaymentRequestReq{
			  PaymentRequestReq{
				  RequestType: "BALANCE",
				  DueDate: dueDate,
				  AutomaticPaymentSource: "NONE",
			  },
		},
		AcceptedPaymentMethods: AcceptedPaymentMethods{
			BankAccount: true,
			BuyNowPayLater: true,
			Card: true,
			SquareGiftCard: true,
		},
		DeliveryMethod: "SHARE_MANUALLY",
		Title: invoiceTitle,
	}
	invoiceJSON, err := json.Marshal(map[string]interface{}{
		"invoice": invoice,
		"idempotency_key": uuid.New().String(),
	})
	if err != nil {
		panic(err)
	}

	body, statusCode := api.createSquareRequest(Params{Method: http.MethodPost, Endpoint: "/invoices", Data: invoiceJSON})
	if statusCode >= 400 {
		var data ErrorResponse
		json.Unmarshal(body, &data)
		log.Println("Error: ", data.Errors)
		return Invoice{}, data.Errors[0].Detail
	} else {
		var mapping InvoiceResponse
		json.Unmarshal(body, &mapping)
		log.Println(mapping)

		return mapping.Invoice, ""
	}
}


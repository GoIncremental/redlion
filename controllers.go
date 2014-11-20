package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/goincremental/redlion/models"
	"github.com/mattbaird/gochimp"
	stripe "github.com/stripe/stripe-go"
	charge "github.com/stripe/stripe-go/charge"
)

func getContextAndRender(template string, w http.ResponseWriter, req *http.Request) {
	c, err := GetRequestContext(req)
	if err != nil {
		log.Printf("Error getting request context: %s\n", err)
	}
	c.Render.HTML(w, http.StatusOK, template, c)
}

//Home renders the default html home page
func Home(w http.ResponseWriter, req *http.Request) {
	getContextAndRender("home", w, req)
}

//Checkout receives an order with a valid stripe token, charges the order via
//stripe and then sends a confirmation email / failure email as required
func Checkout(w http.ResponseWriter, req *http.Request) {
	ctx, err := GetRequestContext(req)

	decoder := json.NewDecoder(req.Body)
	var param models.OrderParams
	err = decoder.Decode(&param)
	if err != nil {
		log.Printf("Error getting order params: %s\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	counter := models.NewCounter(ctx.SiteID, ctx.DB)
	next, err := counter.GetNext("webOrderNumber")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	orderNum := fmt.Sprintf("#%05d", next)
	stripe.Key = os.Getenv("STRIPE_API_KEY")
	ch, err := charge.New(&stripe.ChargeParams{
		Amount:   param.GrandTotal(),
		Currency: "gbp",
		Card:     &stripe.CardParams{Token: param.Token},
		Desc:     orderNum,
	})

	if err != nil {
		stripeErr := err.(*stripe.Error)
		log.Printf("stripe error: %+v\n", stripeErr)
		//TODO: send an admin email
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("STRIPE RESPONSE: %s\n", ch.ID)

	mandrillKey := os.Getenv("MANDRILL_API_KEY")
	mandrillAPI, err := gochimp.NewMandrill(mandrillKey)

	recipients := []gochimp.Recipient{
		gochimp.Recipient{
			Email: param.Email,
			Name:  param.Customer.BillingName,
			Type:  "to",
		},
		gochimp.Recipient{
			Email: "orders@theredlionafterhours.co.uk",
			Name:  "After Hours Order Received",
			Type:  "bcc",
		},
	}
	amt := float64(param.GrandTotal()) / 100
	vars := []gochimp.Var{
		gochimp.Var{
			Name:    "grand_total",
			Content: fmt.Sprintf("£%.2f", amt),
		},
		gochimp.Var{
			Name:    "quantity",
			Content: fmt.Sprintf("%d", param.Quantity),
		},
		gochimp.Var{
			Name:    "invoice_num",
			Content: orderNum,
		},
		gochimp.Var{
			Name:    "customer_name",
			Content: param.Customer.BillingName,
		},
		gochimp.Var{
			Name:    "post",
			Content: strconv.FormatBool(param.Post),
		},
	}

	message := gochimp.Message{
		Subject:         "Your after hours calendar order confirmation " + orderNum,
		FromEmail:       "orders@theredlionafterhours.co.uk",
		FromName:        "The Red Lion",
		To:              recipients,
		GlobalMergeVars: vars,
	}

	sendResponse, err := mandrillAPI.MessageSendTemplate("order-confirmation", []gochimp.Var{}, message, false)
	if err != nil {
		log.Printf("error from mandrill: %s\n", err)
	}

	log.Printf("send response: %+v\n", sendResponse)

	result := &models.OrderResult{
		ID:     orderNum,
		Status: "payment completed",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	j, _ := json.Marshal(result)
	w.Write(j)
}

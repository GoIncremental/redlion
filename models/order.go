package models

type StripeParam struct {
	BillingCity        string `json:"billing_address_city"`
	BillingCountry     string `json:"billing_address_country"`
	BillingCountryCode string `json:"billing_address_country_code"`
	BillingAddress     string `json:"billing_address_line1"`
	BillingState       string `json:"billing_address_state"`
	BillingPostCode    string `json:"billing_address_zip"`
	BillingName        string `json:"billing_name"`
}

type OrderParams struct {
	Quantity uint64      `json:"quantity"`
	Post     bool        `json:"byPost"`
	Token    string      `json:"token"`
	Email    string      `json:"email"`
	Customer StripeParam `json:"args"`
}

func (o *OrderParams) GrandTotal() uint64 {
	result := o.Quantity*1000 + o.Quantity*25 + 20
	if o.Post {
		result += 300
	}
	return result
}

type OrderResult struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

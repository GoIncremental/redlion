package models

import "github.com/goincremental/dal"

type StripeParam struct {
	BillingCity         string `bson:"b_city" json:"billing_address_city"`
	BillingCountry      string `bson:"b_country" json:"billing_address_country"`
	BillingCountryCode  string `bson:"b_country_code" json:"billing_address_country_code"`
	BillingAddress      string `bson:"b_line1" json:"billing_address_line1"`
	BillingState        string `bson:"b_state" json:"billing_address_state"`
	BillingPostCode     string `bson:"b_code" json:"billing_address_zip"`
	BillingName         string `bson:"b_name" json:"billing_name"`
	ShippingCity        string `bson:"s_city" json:"shipping_address_city"`
	ShippingCountry     string `bson:"s_country" json:"shipping_address_country"`
	ShippingCountryCode string `bson:"s_country_code" json:"shipping_address_country_code"`
	ShippingAddress     string `bson:"s_line1" json:"shipping_address_line1"`
	ShippingState       string `bson:"s_state" json:"shipping_address_state"`
	ShippingPostCode    string `bson:"s_code" json:"shipping_address_zip"`
	ShippingName        string `bson:"s_name" json:"shipping_name"`
}

type OrderParams struct {
	ID       dal.ObjectID `bson:"_id" json:"id,omitempty"`
	Quantity uint64       `bson:"quantity" json:"quantity"`
	Post     bool         `bson:"post" json:"post"`
	Token    string       `bson:"token" json:"token"`
	Email    string       `bson:"email" json:"email"`
	Phone    string       `bson:"phone" json:"phone"`
	Customer StripeParam  `bson:",inline" json:"args"`
	Result   OrderResult  `bson:",inline"`
}

func (o *OrderParams) GrandTotal() uint64 {
	result := o.Quantity*1000 + o.Quantity*25 + 25
	if o.Post {
		result += 500
	}
	return result
}

type OrderResult struct {
	ID     string `bson:"reference" json:"id"`
	Status string `bson:"status" json:"status"`
}

type orderParams struct {
	siteID *dal.ObjectID
	col    dal.Collection
}

// OrderParamsModel defines CRUD methods on the OrderParams struct
type OrderParamsModel interface {
	Save(*OrderParams) error
}

func (m *orderParams) Save(item *OrderParams) (err error) {
	if !item.ID.Valid() {
		item.ID = dal.NewObjectID()
	}
	_, err = m.col.SaveID(item.ID, item)
	return
}

// NewOrderParams returns a OrderParamsModel struct pre configured with a
// dal database connection and siteID ready to run the CRUD queries
func NewOrderParams(siteID *dal.ObjectID, db dal.Database) OrderParamsModel {
	col := db.C("orders")
	return &orderParams{
		siteID: siteID,
		col:    col,
	}
}

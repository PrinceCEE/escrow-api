package models

const (
	TransactionTypeProduct = "Product"
	TransactionTypeService = "Service"
	TransactionTypeCrypto  = "Crypto"
)

const (
	TransactionCreatedBySeller = "Seller"
	TransactionCreatedByBuyer  = "Buyer"
)

const (
	TransactionStatusAwaiting        = "Sent-Awaiting"
	TransactionStatusPendingPayment  = "Pending-Payment"
	TransactionStatusPendingDelivery = "Pending-Delivery"
	TransactionStatusCanceled        = "Canceled"
	TransactionStatusCompleted       = "Completed"
)

type ChargeConfiguration struct {
	BuyerCharges  int
	SellerCharges int
}

type ProductDetail struct {
	Name        string
	Quantity    int
	Description string
	Price       int
}

type Transaction struct {
	Status              string              `json:"status" db:"status"`
	Type                string              `json:"type" db:"type"`
	CreatedBy           string              `json:"created_by" db:"created_by"`
	BuyerID             string              `json:"buyer_id" db:"buyer_id"`
	SellerID            string              `json:"seller_id" db:"seller_id"`
	DeliveryDuration    int                 `json:"delivery_duration" db:"delivery_duration"` // in days
	Currency            string              `json:"currency" db:"currency"`
	ChargeConfiguration ChargeConfiguration `json:"charge_configuration" db:"charge_configuration"` // charge configuration in percentage
	ProductDetails      []ProductDetail     `json:"product_details" db:"product_details"`
	TotalAmount         int                 `json:"total_amount" db:"total_amount"`
	TotalCost           int                 `json:"total_cost" db:"total_cost"`
	Charges             int                 `json:"charges" db:"charges"`
	ReceivableAmount    int                 `json:"receivable_amount" db:"receivable_amount"`
	Seller              *Business           `json:"seller" db:"-"`
	Buyer               *User               `json:"buyer" db:"-"`
	Timeline            []*TransactionTimeline
	ModelMixin
}

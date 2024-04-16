package models

const (
	TimelineCreated          = "Transaction Created"
	TimelineApproved         = "Transaction Approved"
	TimelinePaymentSubmitted = "Payment Submitted"
	TimelineDeliveryDone     = "Delivery Done"
	TimelineCompleted        = "Marked As Completed"
	TImelineCanceled         = "Transaction Canceled"
)

type TransactionTimeline struct {
	Name          string       `json:"name" db:"name"`
	TransactionID string       `json:"transaction_id" db:"transaction_id"`
	Transaction   *Transaction `json:"transaction" db:"-"`
	ModelMixin
}

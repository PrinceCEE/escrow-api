package transactions

type createTransactionDto struct {
	Type                string `json:"type" validate:"required,alpha,oneof=Product Service Crypto"`
	CreatedBy           string `json:"created_by" validate:"required,alpha,oneof=Seller Buyer"`
	BuyerID             string `json:"buyer_id,omitempty" validate:"omitempty,uuid"`
	SellerID            string `json:"seller_id,omitempty" validate:"omitempty,uuid"`
	DeliveryDuration    int    `json:"delivery_duration" validate:"required,int"`
	Currency            string `json:"currency" validate:"required,oneof=NGN"`
	ChargeConfiguration struct {
		BuyerCharges  int `json:"buyer_charges" validate:"int,required"`
		SellerCharges int `json:"seller_charges" validate:"int,required"`
	} `json:"charge_configuration" validate:"required"`
	ProductDetails []struct {
		Name        string `json:"name" validate:"required,alphanum"`
		Quantity    int    `json:"quantity" validate:"omitempty,int"`
		Description string `json:"description" validate:"required,alphanum"`
		Price       int    `json:"price" validate:"omitempty,int"`
	} `json:"product_details" validate:"dive"`
}

type updateTransactionDto struct {
	DeliveryDuration    *int    `json:"delivery_duration" validate:"omitempty,int"`
	Currency            *string `json:"currency" validate:"omitempty,oneof=NGN"`
	ChargeConfiguration *struct {
		BuyerCharges  int `json:"buyer_charges" validate:"int,required"`
		SellerCharges int `json:"seller_charges" validate:"int,required"`
	} `json:"charge_configuration" validate:"omitempty"`
	ProductDetails []struct {
		Name        string `json:"name" validate:"required,alphanum"`
		Quantity    int    `json:"quantity" validate:"omitempty,int"`
		Description string `json:"description" validate:"required,alphanum"`
		Price       int    `json:"price" validate:"omitempty,int"`
	} `json:"product_details" validate:"omitempty,dive"`
	Status *string `json:"status,omitempty" validate:"omitempty"`
}

type makePaymentDto struct {
	TransactionID string `json:"transaction_id" validate:"required,uuid"`
	IsUseWallet   bool   `json:"is_use_wallet" validate:"required,bool"`
}

type getTransactionsQueryDto struct {
	Page      int    `json:"page" validate:"number,min=1"`
	PageSize  int    `json:"page_size" validate:"number,min=1,max=100"`
	Status    string `json:"status"`
	Type      string `json:"type"`
	CreatedBy string `json:"created_by"`
	BuyerId   string `json:"buyer_id"`
	SellerId  string `json:"seller_id"`
}

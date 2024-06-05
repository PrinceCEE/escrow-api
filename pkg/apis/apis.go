package apis

import "github.com/princecee/escrow-api/pkg/apis/paystack"

type IAPIs interface {
	GetPaystack() paystack.IPaystack
}

type apis struct{}

func NewAPIs() *apis {
	return &apis{}
}

func (a *apis) GetPaystack() paystack.IPaystack {
	return paystack.NewPaystackAPI()
}

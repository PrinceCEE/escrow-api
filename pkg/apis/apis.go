package apis

import "github.com/Bupher-Co/bupher-api/pkg/apis/paystack"

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

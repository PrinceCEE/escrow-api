package push

type Sms struct {
	Phone   string
	Message string
}

func SendSMS(data *Sms) {}

package models

type EventEnvironment string
type EventType string

const (
	AppEnvironment  EventEnvironment = "app_environment"
	PushEnvironment EventEnvironment = "push_environment"
	JobEnvironment  EventEnvironment = "job_environment"
)

const (
	SmsEventType   EventType = "sms"
	EmailEventType EventType = "email"
)

type Event struct {
	Data              any              `json:"data"`
	OriginEnvironment EventEnvironment `json:"origin_environment"`
	TargetEnvironment EventEnvironment `json:"target_environment"`
	EventType         EventType        `json:"event_type"`
	ModelMixin
}

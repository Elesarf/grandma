package main

// InternalMessage внутреннее представление сообщения
type InternalMessage struct {
	chatID      int64
	messageID   int
	userName    string
	messageText string
}

// InternalLocation internal location entity
type InternalLocation struct {
	mainLocation     string
	optionalLocation string
}

// InternalUser internal user entity
type InternalUser struct {
	userName  string
	location  InternalLocation
	timestamp uint64
}

// NewInternalUser constructor
func NewInternalUser(name string, loc InternalLocation, time uint64) *InternalUser {
	return &InternalUser{name, loc, time}
}

package telegramBot

func NewRecipient(recipient string) *Recipient {
	return &Recipient{recipient: recipient}
}

type Recipient struct {
	recipient string
}

func (r *Recipient) Recipient() string {
	return r.recipient
}

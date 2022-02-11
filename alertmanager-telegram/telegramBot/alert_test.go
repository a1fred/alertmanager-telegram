package telegramBot_test

import (
	"testing"
	"time"

	"github.com/a1fred/alertmanager-telegram/alertmanager-telegram/telegramBot"
	"github.com/jaswdr/faker"

	"github.com/prometheus/alertmanager/notify/webhook"
	"github.com/prometheus/alertmanager/template"
	"github.com/stretchr/testify/assert"
)

var fake faker.Faker = faker.New()

func FakeAlert() template.Alert {
	startAt := fake.Time().Time(time.Now())

	return template.Alert{
		Status:       fake.RandomStringElement([]string{"firing", "resolved"}),
		Labels:       map[string]string{fake.Company().JobTitle(): fake.Company().JobTitle()},
		Annotations:  map[string]string{fake.Company().JobTitle(): fake.Company().JobTitle()},
		StartsAt:     startAt,
		EndsAt:       fake.Time().Time(startAt),
		GeneratorURL: "http://localhost:9093",
		Fingerprint:  fake.UUID().V4(),
	}

}

func FakeAlertSlice(size int) template.Alerts {
	result := make(template.Alerts, 0)
	for i := 0; i <= size; i++ {
		result = append(result, FakeAlert())
	}

	return result
}

func FakeWebhookMessage() webhook.Message {
	status := fake.RandomStringElement([]string{"firing", "resolved"})

	return webhook.Message{
		Version:         "4",
		GroupKey:        "dummyGroup",
		TruncatedAlerts: 5,
		Data: &template.Data{
			Receiver: "telegram",
			Status:   status,
			Alerts:   FakeAlertSlice(fake.IntBetween(1, 10)),

			GroupLabels:       map[string]string{},
			CommonLabels:      map[string]string{},
			CommonAnnotations: map[string]string{},
			ExternalURL:       "http://localhost:9093",
		},
	}
}

func TestFormatAlertHtml(t *testing.T) {
	msg, err := telegramBot.FormatAlertHtml(FakeWebhookMessage(), time.UTC)
	assert.NoError(t, err)
	assert.NotZero(t, msg)
}

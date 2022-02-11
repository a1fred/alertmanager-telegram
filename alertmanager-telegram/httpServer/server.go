package httpServer

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/prometheus/alertmanager/notify/webhook"
	"github.com/prometheus/client_golang/prometheus"
)

func NewHttpServeMux(httpLogger *log.Logger, alertmanagerMessages chan<- webhook.Message, alertsReceivedCounter prometheus.Counter) *http.ServeMux {
	m := http.NewServeMux()

	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if r.Body == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var message webhook.Message

		err := json.NewDecoder(r.Body).Decode(&message)
		if err != nil {
			httpLogger.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		alertsReceivedCounter.Add(float64(len(message.Alerts)))
		alertmanagerMessages <- message
		w.WriteHeader(http.StatusOK)
	})
	m.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return m
}

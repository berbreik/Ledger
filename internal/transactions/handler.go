package transactions

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"net/http"
)

type Handler struct {
	MQChannel *amqp.Channel
	QueueName string
}

func (h *Handler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var tx Transaction
	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		http.Error(w, "Invalid Payload", http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(tx)
	if err != nil {
		http.Error(w, "Failed to encode transaction", http.StatusInternalServerError)
		return
	}

	err = h.MQChannel.Publish(
		"",
		h.QueueName,
		false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		},
	)
	if err != nil {
		http.Error(w, "Failed to queue transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Transaction queued successfully "))
}

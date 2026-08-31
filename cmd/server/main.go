package main

import (
	"fmt"
	"net/http"

	"broker/internal/broker"
)

func main() {
	b := broker.NewBroker()

	http.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		topic := r.URL.Query().Get("topic")
		message := r.URL.Query().Get("message")

		if topic == "" || message == "" {
			http.Error(w, "topic and message are required", http.StatusBadRequest)
			return
		}

		b.Publish(topic, message)
		fmt.Fprintln(w, "published")
	})

	http.HandleFunc("/consume", func(w http.ResponseWriter, r *http.Request) {
		topic := r.URL.Query().Get("topic")
		if topic == "" {
			http.Error(w, "topic is required", http.StatusBadRequest)
			return
		}

		msg, ok := b.Consume(topic)
		if !ok {
			http.Error(w, "no messages", http.StatusNoContent)
			return
		}
		fmt.Fprintln(w, msg)
	})

	http.HandleFunc("/count", func(w http.ResponseWriter, r *http.Request) {
		topic := r.URL.Query().Get("topic")
		fmt.Fprintln(w, b.Count(topic))
	})
	fmt.Println("listening on :8080")
	http.ListenAndServe(":8080", nil)

}

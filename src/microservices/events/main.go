package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/IBM/sarama"
)

// (KAFKA_CREATE_TOPICS в docker-compose).
var topics = map[string]string{
	"movie":   "movie-events",
	"user":    "user-events",
	"payment": "payment-events",
}

// Event — конверт события по спецификации API.
type Event struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	Payload   map[string]interface{} `json:"payload"`
}

// EventResponse — ответ на успешную публикацию.
type EventResponse struct {
	Status    string `json:"status"`
	Partition int32  `json:"partition"`
	Offset    int64  `json:"offset"`
	Event     Event  `json:"event"`
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func writeJSON(w http.ResponseWriter, code int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(body)
}

// waitProducer ждёт готовности Kafka
func waitProducer(brokers []string, cfg *sarama.Config) sarama.SyncProducer {
	for i := 1; i <= 30; i++ {
		p, err := sarama.NewSyncProducer(brokers, cfg)
		if err == nil {
			log.Printf("producer: подключен к Kafka %v", brokers)
			return p
		}
		log.Printf("producer: Kafka недоступна (попытка %d/30): %v", i, err)
		time.Sleep(2 * time.Second)
	}
	log.Fatal("producer: не удалось подключиться к Kafka")
	return nil
}

func waitConsumer(brokers []string, cfg *sarama.Config) sarama.Consumer {
	for i := 1; i <= 30; i++ {
		c, err := sarama.NewConsumer(brokers, cfg)
		if err == nil {
			log.Printf("consumer: подключен к Kafka %v", brokers)
			return c
		}
		log.Printf("consumer: Kafka недоступна (попытка %d/30): %v", i, err)
		time.Sleep(2 * time.Second)
	}
	log.Fatal("consumer: не удалось подключиться к Kafka")
	return nil
}

// consume читает топик и пишет всё прочитанное в лог сервиса.
func consume(consumer sarama.Consumer, topic string) {
	for {
		pc, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
		if err != nil {
			log.Printf("consumer: не удалось подписаться на %s: %v", topic, err)
			time.Sleep(3 * time.Second)
			continue
		}
		log.Printf("consumer: подписан на топик %s", topic)

		for msg := range pc.Messages() {
			log.Printf("consumer: топик=%s партиция=%d offset=%d событие=%s",
				msg.Topic, msg.Partition, msg.Offset, string(msg.Value))
		}

		pc.Close()
		log.Printf("consumer: соединение с %s закрыто, переподключаюсь", topic)
	}
}

// handleEvent принимает событие, кладёт его в свой топик и возвращает
// партицию и offset, куда оно легло.
func handleEvent(producer sarama.SyncProducer, eventType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method Not Allowed"})
			return
		}

		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "некорректный JSON: " + err.Error()})
			return
		}

		event := Event{
			ID:        fmt.Sprintf("%s-%d", eventType, time.Now().UnixNano()),
			Type:      eventType,
			Timestamp: time.Now().UTC(),
			Payload:   payload,
		}

		data, err := json.Marshal(event)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
			Topic: topics[eventType],
			Value: sarama.ByteEncoder(data),
		})
		if err != nil {
			log.Printf("producer: ошибка отправки в %s: %v", topics[eventType], err)
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		log.Printf("producer: топик=%s партиция=%d offset=%d событие=%s",
			topics[eventType], partition, offset, event.ID)

		writeJSON(w, http.StatusCreated, EventResponse{
			Status:    "success",
			Partition: partition,
			Offset:    offset,
			Event:     event,
		})
	}
}

func main() {
	port := getenv("PORT", "8082")
	brokers := strings.Split(getenv("KAFKA_BROKERS", "kafka:9092"), ",")

	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_7_0_0        // версия брокера из docker-compose
	cfg.Producer.Return.Successes = true // обязательно для SyncProducer
	cfg.Producer.RequiredAcks = sarama.WaitForLocal

	producer := waitProducer(brokers, cfg)
	defer producer.Close()

	consumer := waitConsumer(brokers, cfg)
	defer consumer.Close()

	// Сервис сам читает то, что сам же и написал.
	for _, topic := range topics {
		go consume(consumer, topic)
	}

	http.HandleFunc("/api/events/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]bool{"status": true})
	})
	http.HandleFunc("/api/events/movie", handleEvent(producer, "movie"))
	http.HandleFunc("/api/events/user", handleEvent(producer, "user"))
	http.HandleFunc("/api/events/payment", handleEvent(producer, "payment"))

	log.Printf("events-service: слушаю порт %s, брокеры %v", port, brokers)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("events-service: сервер остановлен: %v", err)
	}
}

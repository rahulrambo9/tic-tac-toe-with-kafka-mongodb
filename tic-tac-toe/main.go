package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/segmentio/kafka-go"
)

type GameResult struct {
    Player1Name string `json:"player1_name"`
    Player2Name string `json:"player2_name"`
    WinnerName  string `json:"winner_name"`
}

// Kafka producer to send game results
func publishToKafka(result GameResult) error {
    kafkaURL := "35.237.164.110:9092"
    topic := "game-results"
    
    w := kafka.Writer{
        Addr:     kafka.TCP(kafkaURL),
        Topic:    topic,
        Balancer: &kafka.LeastBytes{},
    }

    message, err := json.Marshal(result)
    if err != nil {
        return err
    }

    err = w.WriteMessages(context.Background(), kafka.Message{
        Key:   []byte(result.WinnerName),
        Value: message,
    })

    if err != nil {
        return fmt.Errorf("failed to publish to Kafka: %v", err)
    }

    fmt.Println("Game result published to Kafka:", result)
    return nil
}

// Handle saving game result and sending to Kafka
func saveGameResult(w http.ResponseWriter, r *http.Request) {
    var result GameResult
    err := json.NewDecoder(r.Body).Decode(&result)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    err = publishToKafka(result)
    if err != nil {
        http.Error(w, "Failed to publish game result to Kafka", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"message": "Game result sent to Kafka"}`))
}

func main() {
    fs := http.FileServer(http.Dir("./static"))
    http.Handle("/", fs)

    http.HandleFunc("/save-result", saveGameResult)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    fmt.Printf("Server running on http://localhost:%s...\n", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}

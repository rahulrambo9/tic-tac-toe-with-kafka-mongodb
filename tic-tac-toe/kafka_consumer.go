package main

import (
    "context"
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "time"

    _ "github.com/lib/pq"  // Anonymous import for the PostgreSQL driver
    "github.com/segmentio/kafka-go"
)

type GameResult struct {
    Player1Name string `json:"player1_name"`
    Player2Name string `json:"player2_name"`
    WinnerName  string `json:"winner_name"`
}

// Connect to PostgreSQL with user, password, and dbname
func connectDB() (*sql.DB, error) {
    connStr := "user=postgres password=UTtf7mWwiy host=host.docker.internal port=5432 dbname=tic_tac_toe_db sslmode=disable"
    return sql.Open("postgres", connStr)
}

// Kafka consumer to read game results and store them in PostgreSQL
func consumeFromKafka() {
    kafkaURL := "35.237.164.110:9092"
    topic := "game-results"

    r := kafka.NewReader(kafka.ReaderConfig{
        Brokers:   []string{kafkaURL},
        Topic:     topic,
        Partition: 0,
        MinBytes:  10e3,
        MaxBytes:  10e6,
    })

    db, err := connectDB()
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()

    fmt.Println("Listening to Kafka topic:", topic)

    for {
        msg, err := r.ReadMessage(context.Background())
        if err != nil {
            log.Fatal("Error reading message from Kafka:", err)
        }

        var result GameResult
        err = json.Unmarshal(msg.Value, &result)
        if err != nil {
            log.Println("Error unmarshaling Kafka message:", err)
            continue
        }

        // Store result in PostgreSQL
        err = storeGameResult(db, result)
        if err != nil {
            log.Println("Error storing result in database:", err)
            continue
        }

        fmt.Println("Game result stored in DB:", result)
    }
}

// Store game result in PostgreSQL
func storeGameResult(db *sql.DB, result GameResult) error {
    query := `INSERT INTO games (player1_name, player2_name, winner_name, game_date) VALUES ($1, $2, $3, $4)`
    _, err := db.Exec(query, result.Player1Name, result.Player2Name, result.WinnerName, time.Now())
    return err
}

func main() {
    consumeFromKafka()
}

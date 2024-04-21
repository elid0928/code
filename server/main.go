package main

import (
	"bufio"
	"bytes"
	"context"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
)

const (
	defaultPort    = "8080"
	defaultTimeout = 30 * time.Second
)

const (
	apiKey = "xxx"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	addr := ":" + port
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to bind to %s: %v", addr, err)
	}
	defer l.Close()

	log.Printf("Listening on %s", addr)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		go handleConnection(conn, apiKey)
	}
}

func handleConnection(conn net.Conn, apiKey string) {
	defer conn.Close()

	client := conn.RemoteAddr().String()
	log.Printf("New client connected: %s", client)

	c := openai.NewClient(apiKey)
	// c.Timeout = defaultTimeout

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Text()
		message = strings.TrimSpace(message)

		if message == "/quit" {
			log.Printf("Client left: %s", client)
			return
		}

		respMsg, err := askToGPT(c, message)
		if err != nil {
			log.Printf("Failed to send message to GPT-3: %v", err)
			continue
		}

		var buf bytes.Buffer
		buf.WriteString("(echo) gpt: ")
		buf.WriteString(respMsg)
		buf.WriteString("%")
		reply := buf.Bytes()

		if _, err := conn.Write(reply); err != nil {
			log.Printf("Failed to write reply to the client: %v", err)
			return
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Failed to read from client: %v", err)
	}
}

func askToGPT(client *openai.Client, prompt string) (string, error) {
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

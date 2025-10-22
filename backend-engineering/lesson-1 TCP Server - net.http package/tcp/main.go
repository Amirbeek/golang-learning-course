package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

// Shared key-value store (thread-safe)
var (
	mu sync.RWMutex
	db = make(map[string]string)
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("✅ TCP Server started on port :8080")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Accept error:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	writeLine(conn, "Welcome to Go TCP Demo Server!")
	writeLine(conn, "Type HELP to see all commands.")

	reader := bufio.NewReader(conn)
	for {
		writeRaw(conn, "> ")

		line, err := reader.ReadString('\n')
		if err != nil {
			return // client disconnected
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		cmd := strings.ToUpper(parts[0])

		switch cmd {
		case "PING":
			writeLine(conn, "PONG")

		case "ECHO":
			if len(parts) < 2 {
				writeLine(conn, "ERR usage: ECHO <message>")
				continue
			}
			msg := strings.Join(parts[1:], " ")
			writeLine(conn, msg)

		case "SET":
			if len(parts) < 3 {
				writeLine(conn, "ERR usage: SET <key> <value>")
				continue
			}
			key := parts[1]
			value := strings.Join(parts[2:], " ")
			mu.Lock()
			db[key] = value
			mu.Unlock()
			writeLine(conn, "OK")

		case "GET":
			if len(parts) != 2 {
				writeLine(conn, "ERR usage: GET <key>")
				continue
			}
			key := parts[1]
			mu.RLock()
			value, ok := db[key]
			mu.RUnlock()
			if ok {
				writeLine(conn, value)
			} else {
				writeLine(conn, "(nil)")
			}

		case "DEL":
			if len(parts) < 2 {
				writeLine(conn, "ERR usage: DEL <key> [key2 ...]")
				continue
			}
			count := 0
			mu.Lock()
			for _, key := range parts[1:] {
				if _, ok := db[key]; ok {
					delete(db, key)
					count++
				}
			}
			mu.Unlock()
			writeLine(conn, fmt.Sprintf("Deleted %d key(s)", count))

		case "TIME":
			writeLine(conn, time.Now().Format(time.RFC3339))

		case "HELP":
			writeLine(conn, "Available commands:")
			writeLine(conn, "  PING                -> responds with PONG")
			writeLine(conn, "  ECHO <text>         -> echoes your text")
			writeLine(conn, "  SET <key> <value>   -> saves a key-value pair")
			writeLine(conn, "  GET <key>           -> returns stored value")
			writeLine(conn, "  DEL <key> [key2..]  -> deletes one or more keys")
			writeLine(conn, "  TIME                -> shows current server time")
			writeLine(conn, "  QUIT                -> closes the connection")

		case "QUIT", "EXIT":
			writeLine(conn, "Goodbye!")
			return

		default:
			writeLine(conn, "ERR unknown command. Type HELP.")
		}
	}
}

// Helper to write a line with newline
func writeLine(conn net.Conn, s string) {
	writeRaw(conn, s+"\r\n")
}

// Helper to write raw text
func writeRaw(conn net.Conn, s string) {
	_, _ = conn.Write([]byte(s))
}

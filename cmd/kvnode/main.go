package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "time"

    "memora/internal/storage"

)

func main() {
    store, err := storage.Open("data.wal")
    if err != nil {
        panic(err)
    }
    defer store.Close()

    fmt.Println("KV Store started")
    fmt.Println("Commands: PUT key value | GET key | DEL key | EXIT")

    scanner := bufio.NewScanner(os.Stdin)
    for {
        fmt.Print("> ")
        if !scanner.Scan() {
            break
        }

        input := strings.TrimSpace(scanner.Text())
        parts := strings.SplitN(input, " ", 3)

        if len(parts) == 0 {
            continue
        }

        switch strings.ToUpper(parts[0]) {
        case "PUT":
            if len(parts) != 3 {
                fmt.Println("Usage: PUT key value")
                continue
            }
            err := store.Put(parts[1], []byte(parts[2]), 0)
            if err != nil {
                fmt.Println("ERR:", err)
            } else {
                fmt.Println("OK")
            }

        case "GET":
            if len(parts) != 2 {
                fmt.Println("Usage: GET key")
                continue
            }
            val, ok := store.Get(parts[1])
            if !ok {
                fmt.Println("(nil)")
            } else {
                fmt.Println(string(val))
            }

        case "DEL":
            if len(parts) != 2 {
                fmt.Println("Usage: DEL key")
                continue
            }
            err := store.Delete(parts[1])
            if err != nil {
                fmt.Println("ERR:", err)
            } else {
                fmt.Println("OK")
            }

        case "EXIT":
            fmt.Println("Shutting down")
            time.Sleep(100 * time.Millisecond) // allow fsync
            return

        default:
            fmt.Println("Unknown command")
        }
    }
}

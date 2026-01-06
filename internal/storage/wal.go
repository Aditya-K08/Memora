package storage

import (
    "bufio"
    "encoding/binary"
    "os"
    "sync"
)
const (
    opPut byte = 1
    opDel byte = 2
)

type WAL struct {
    mu     sync.Mutex
    file   *os.File
    writer *bufio.Writer
}

func OpenWAL(path string) (*WAL, error) {
    f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
    if err != nil {
        return nil, err
    }
    return &WAL{
        file:   f,
        writer: bufio.NewWriter(f),
    }, nil
}

func (w *WAL) AppendPut(key string, value []byte, expiry int64) error {
    w.mu.Lock()
    defer w.mu.Unlock()

    if err := w.writer.WriteByte(opPut); err != nil {
        return err
    }
    if err := binary.Write(w.writer, binary.BigEndian, uint32(len(key))); err != nil {
        return err
    }
    if _, err := w.writer.WriteString(key); err != nil {
        return err
    }
    if err := binary.Write(w.writer, binary.BigEndian, uint32(len(value))); err != nil {
        return err
    }
    if _, err := w.writer.Write(value); err != nil {
        return err
    }
    if err := binary.Write(w.writer, binary.BigEndian, expiry); err != nil {
        return err
    }

    if err := w.writer.Flush(); err != nil {
        return err
    }
    return w.file.Sync()
}

func (w *WAL) AppendDelete(key string) error {
    w.mu.Lock()
    defer w.mu.Unlock()

    if err := w.writer.WriteByte(opDel); err != nil {
        return err
    }
    if err := binary.Write(w.writer, binary.BigEndian, uint32(len(key))); err != nil {
        return err
    }
    if _, err := w.writer.WriteString(key); err != nil {
        return err
    }

    if err := w.writer.Flush(); err != nil {
        return err
    }
    return w.file.Sync()
}

func (w *WAL) Close() error {
    w.mu.Lock()
    defer w.mu.Unlock()
    w.writer.Flush()
    return w.file.Close()
}

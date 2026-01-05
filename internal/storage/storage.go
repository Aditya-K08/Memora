package storage

import (
    "bufio"
    "encoding/binary"
    "io"
    "os"
    "time"
)

type KVStore struct {
    mem *MemTable
    wal *WAL
}

func Open(path string) (*KVStore, error) {
    mem := NewMemTable()

    f, err := os.Open(path)
    if err == nil {
        reader := bufio.NewReader(f)
        for {
            op, err := reader.ReadByte()
            if err == io.EOF {
                break
            }
            if err != nil {
                return nil, err
            }

            var keyLen uint32
            binary.Read(reader, binary.BigEndian, &keyLen)
            key := make([]byte, keyLen)
            reader.Read(key)

            if op == opPut {
                var valLen uint32
                binary.Read(reader, binary.BigEndian, &valLen)
                val := make([]byte, valLen)
                reader.Read(val)

                var exp int64
                binary.Read(reader, binary.BigEndian, &exp)

                mem.data[string(key)] = Value{Data: val, Expiry: exp}
            } else {
                delete(mem.data, string(key))
            }
        }
        f.Close()
    }

    wal, err := OpenWAL(path)
    if err != nil {
        return nil, err
    }

    return &KVStore{mem: mem, wal: wal}, nil
}

func (k *KVStore) Get(key string) ([]byte, bool) {
    return k.mem.Get(key)
}

func (k *KVStore) Put(key string, val []byte, ttl time.Duration) error {
    var exp int64
    if ttl > 0 {
        exp = time.Now().Add(ttl).UnixNano()
    }
    if err := k.wal.AppendPut(key, val, exp); err != nil {
        return err
    }
    k.mem.Put(key, val, ttl)
    return nil
}

func (k *KVStore) Delete(key string) error {
    if err := k.wal.AppendDelete(key); err != nil {
        return err
    }
    k.mem.Delete(key)
    return nil
}

func (k *KVStore) Close() error {
    return k.wal.Close()
}

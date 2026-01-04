package utils

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"os"
	"runtime"
	"time"
)

func GenerateQueryID() string {
	now := time.Now().UnixNano()
	pid := int32(os.Getpid())
	tid := int32(runtime.NumGoroutine())

	buf := make([]byte, 16)
	binary.LittleEndian.PutUint64(buf[:8], uint64(now))
	binary.LittleEndian.PutUint32(buf[8:12], uint32(pid))
	binary.LittleEndian.PutUint32(buf[12:16], uint32(tid))

	hash := sha256.Sum256(buf)
	return fmt.Sprintf("%x", hash[:12])
}

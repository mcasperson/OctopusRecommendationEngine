package wait

import (
	"fmt"
	"time"
)

func WaitForResource(callback func() error, timeout time.Duration) error {
	start := time.Now()
	for {
		err := callback()
		if err == nil {
			return nil
		}
		time.Sleep(time.Second)
		now := time.Now()
		if now.Sub(start) > timeout {
			break
		}
	}

	return fmt.Errorf("server did not reply after %v", timeout)
}

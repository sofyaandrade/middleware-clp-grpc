package modbusSlave

import (
	"context"
	"strings"
	"time"
)

func ConnectWithRetry(ctx context.Context, client slaveConnector) error {
	client.Close()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := client.Connect()
			if err == nil {
				return nil
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(2 * time.Second):
			}
		}
	}
}

func isConnectionError(err error) bool {
	if err == nil {
		return false
	}
	if strings.Contains(err.Error(), "connection reset") ||
		strings.Contains(err.Error(), "broken pipe") ||
		strings.Contains(err.Error(), "connection refused") ||
		strings.Contains(err.Error(), "i/o timeout") ||
		strings.Contains(err.Error(), "EOF") ||
		strings.Contains(err.Error(), "forcibly closed") ||
		strings.Contains(err.Error(), "closed by the remote host") ||
		strings.Contains(err.Error(), "wsasend") ||
		strings.Contains(err.Error(), "use of closed network connection") {
		return true
	}

	return false
}

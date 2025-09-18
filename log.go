package ngtelgcp

import (
	"context"
)

func GetLogArgs(ctx context.Context) []any {
	tracePath := GetTracePath(ctx)

	if tracePath == "" {
		return nil
	}

	return []any{"trace", tracePath}
}

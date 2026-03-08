package asynqworker

import (
	"encoding/json"
	"fmt"
	"strconv"

	"sociomile-be/internal/domain/model"

	"github.com/hibiken/asynq"
)

func decodeEvent(task *asynq.Task) (model.DomainEvent, error) {
	var event model.DomainEvent
	if err := json.Unmarshal(task.Payload(), &event); err != nil {
		return model.DomainEvent{}, fmt.Errorf("decode event payload: %w", err)
	}

	return event, nil
}

func getInt64(payload map[string]any, key string) (int64, error) {
	v, ok := payload[key]
	if !ok {
		return 0, fmt.Errorf("missing payload key: %s", key)
	}

	switch t := v.(type) {
	case float64:
		return int64(t), nil
	case int64:
		return t, nil
	case int:
		return int64(t), nil
	case json.Number:
		return t.Int64()
	case string:
		val, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid int payload %s: %w", key, err)
		}
		return val, nil
	default:
		return 0, fmt.Errorf("unsupported payload type for %s", key)
	}
}

func getString(payload map[string]any, key, fallback string) string {
	v, ok := payload[key]
	if !ok {
		return fallback
	}

	s, ok := v.(string)
	if !ok || s == "" {
		return fallback
	}

	return s
}

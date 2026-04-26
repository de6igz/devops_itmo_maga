package httpapi

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	defaultCPULoadDuration = 500 * time.Millisecond
	maxCPULoadDuration     = 5 * time.Second
)

var cpuLoadSlots = make(chan struct{}, 1)

func (h *Handler) cpuLoad(c echo.Context) error {
	duration := defaultCPULoadDuration

	if rawDuration := c.QueryParam("duration"); rawDuration != "" {
		parsedDuration, err := time.ParseDuration(rawDuration)
		if err != nil || parsedDuration <= 0 {
			return echo.NewHTTPError(
				http.StatusBadRequest,
				"duration должен быть положительным значением, например 500ms или 2s",
			)
		}

		duration = parsedDuration
	}

	if duration > maxCPULoadDuration {
		duration = maxCPULoadDuration
	}

	startedAt := time.Now()
	release, err := acquireCPULoadSlot(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "нагрузочный тест отменен клиентом")
	}
	defer release()

	loadStartedAt := time.Now()
	deadline := loadStartedAt.Add(duration)

	workers := 1

	var wg sync.WaitGroup
	wg.Add(workers)

	counters := make([]uint64, workers)
	checksums := make([]uint64, workers)

	for i := 0; i < workers; i++ {
		go func(workerID int) {
			defer wg.Done()

			var counter uint64
			var payload [32]byte

			binary.LittleEndian.PutUint64(payload[:8], uint64(workerID))

			for time.Now().Before(deadline) {
				binary.LittleEndian.PutUint64(payload[:8], counter)
				payload = sha256.Sum256(payload[:])
				counter++
			}

			counters[workerID] = counter
			checksums[workerID] = binary.LittleEndian.Uint64(payload[:8])
		}(i)
	}

	wg.Wait()

	var totalIterations uint64
	var checksum uint64

	for i := 0; i < workers; i++ {
		totalIterations += counters[i]
		checksum ^= checksums[i]
	}

	return c.JSON(http.StatusOK, map[string]any{
		"status":        "ok",
		"workers":       workers,
		"requestedTime": duration.String(),
		"actualTime":    time.Since(startedAt).String(),
		"queueTime":     loadStartedAt.Sub(startedAt).String(),
		"iterations":    totalIterations,
		"checksum":      checksum,
	})
}

func acquireCPULoadSlot(ctx context.Context) (func(), error) {
	select {
	case cpuLoadSlots <- struct{}{}:
		return func() { <-cpuLoadSlots }, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

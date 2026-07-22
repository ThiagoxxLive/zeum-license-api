// Package ratelimit implementa um limitador por janela deslizante, em memória,
// equivalente ao RateLimitSubscriber (Symfony RateLimiter) do zeum-admin-api.
package ratelimit

import (
	"sync"
	"time"
)

type SlidingWindow struct {
	limit    int
	interval time.Duration

	mu  sync.Mutex
	log map[string][]time.Time
}

func NewSlidingWindow(limit int, interval time.Duration) *SlidingWindow {

	return &SlidingWindow{
		limit:    limit,
		interval: interval,
		log:      make(map[string][]time.Time),
	}
}

// Allow registra uma requisição para key (normalmente o IP do cliente) e reporta
// se ela deve ser aceita dentro da janela deslizante, junto com o tempo de espera
// sugerido quando rejeitada.
func (w *SlidingWindow) Allow(key string) (allowed bool, retryAfter time.Duration) {

	now := time.Now()
	windowStart := now.Add(-w.interval)

	w.mu.Lock()
	defer w.mu.Unlock()

	timestamps := w.log[key]
	kept := timestamps[:0]

	for _, t := range timestamps {
		if t.After(windowStart) {
			kept = append(kept, t)
		}
	}

	if len(kept) >= w.limit {
		w.log[key] = kept
		return false, kept[0].Add(w.interval).Sub(now)
	}

	w.log[key] = append(kept, now)

	return true, 0
}

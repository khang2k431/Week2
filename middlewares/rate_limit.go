package middlewares

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	clients = make(map[string]*client)
	mu      sync.Mutex
)

func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()
	c, ok := clients[ip]
	if !ok {
		lim := rate.NewLimiter(5, 10) // 5 req/sec
		clients[ip] = &client{limiter: lim, lastSeen: time.Now()}
		return lim
	}
	c.lastSeen = time.Now()
	return c.limiter
}

func cleanupClients() {
	for {
		time.Sleep(time.Minute)
		mu.Lock()
		for ip, c := range clients {
			if time.Since(c.lastSeen) > 3*time.Minute {
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}
}

func RateLimitMiddleware() gin.HandlerFunc {
	go cleanupClients()
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if net.ParseIP(ip) == nil {
			ip = "unknown"
		}
		lim := getLimiter(ip)
		if !lim.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}
		c.Next()
	}
}

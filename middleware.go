package turispro_user

import (
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthenticatedUser struct {
	ID           string
	Email        string
	TourOperator string
	Level        UserLevel
}

func (u *AuthenticatedUser) IsAdmin() bool {
	return u.Level == Admin
}

func (u *AuthenticatedUser) HasRole(level UserLevel) bool {
	if u.IsAdmin() {
		return true
	}
	return u.Level == level
}

const userContextKey = "user"

func InjectUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		levelStr := c.GetHeader("X-User-Level")
		levelInt, err := strconv.Atoi(levelStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "invalid or missing level header"})
			return
		}

		if levelInt < int(Admin) || levelInt > int(Guia) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "invalid user level"})
			return
		}

		user := &AuthenticatedUser{
			ID:           c.GetHeader("X-User-ID"),
			Email:        c.GetHeader("X-User-Email"),
			TourOperator: c.GetHeader("X-User-TourOperator"),
			Level:        UserLevel(levelInt),
		}

		if user.ID == "" || user.Email == "" || user.TourOperator == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "missing user information"})
			return
		}

		c.Set(userContextKey, user)
		c.Next()
	}
}

func GetUser(c *gin.Context) *AuthenticatedUser {
	if u, exists := c.Get(userContextKey); exists {
		if user, ok := u.(*AuthenticatedUser); ok {
			return user
		}
	}
	return nil
}

func IsInternalCall(c *gin.Context) bool {
	clientIP := c.ClientIP()
	internalHeader := strings.ToLower(c.GetHeader("X-Internal-Call"))

	if isPrivateIP(clientIP) && internalHeader == "true" {
		return true
	}

	return false
}

func isPrivateIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	privateCIDRs := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	}

	for _, cidr := range privateCIDRs {
		if _, subnet, err := net.ParseCIDR(cidr); err == nil && subnet.Contains(ip) {
			return true
		}
	}

	return false
}

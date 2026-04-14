package middleware

import "github.com/gin-gonic/gin"

// SecurityHeaders adds common HTTP security headers to every response.
// This complements HTTPS enforcement done at the reverse-proxy level (Nginx/Caddy).
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent MIME sniffing
		c.Header("X-Content-Type-Options", "nosniff")
		// Prevent clickjacking
		c.Header("X-Frame-Options", "DENY")
		// Enable browser XSS filtering (legacy browsers)
		c.Header("X-XSS-Protection", "1; mode=block")
		// Only send referrer on same-origin requests
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		// Disable browser features that APIs don't need
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		// Strict-Transport-Security: force HTTPS for 1 year including sub-domains
		// c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Next()
	}
}

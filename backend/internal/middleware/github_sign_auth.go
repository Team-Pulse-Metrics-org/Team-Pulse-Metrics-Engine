package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func VerifyGithubSignature(secret string, env string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if env == "development" && secret == "" {
			log.Warn().Msg("Webhook Signature Verification skipped in development mode")
			c.Next()
			return
		}
		signature := c.GetHeader("X-Hub-Signature-256")
		if signature == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid payload signature header"})
			return
		}

		sigParts := strings.SplitN(signature, "=", 2)
		if len(sigParts) != 2 || sigParts[0] != "sha256" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid signature format"})
			return
		}

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write(body)
		expectedMAC := hex.EncodeToString(mac.Sum(nil))

		if !hmac.Equal([]byte(sigParts[1]), []byte(expectedMAC)) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid payload signature"})
			return
		}
		c.Next()
	}
}

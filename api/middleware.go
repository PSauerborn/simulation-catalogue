package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// AuthMiddleware validates API key authentication for admin endpoints.
// It checks the X-API-Key header, verifies the key exists and is not expired
// or revoked, and sets the admin user in the context on success.
func AuthMiddleware(db Persistence) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.Request.Header.Get("X-API-Key")
		ipAddress := c.ClientIP()

		if apiKey == "" {
			log.WithFields(log.Fields{
				"ip_address": ipAddress,
			}).Warn("received request without API key")
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "API key is required"})
			return
		}

		apiKeyObj, err := db.GetAPIKey(apiKey)
		var errKeyNotFound ErrAPIKeyNotFound
		if err != nil && errors.As(err, &errKeyNotFound) {
			log.WithFields(log.Fields{
				"ip_address": ipAddress,
			}).Warn("received request with unknown API key")
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid API key"})
			return

		} else if err != nil {
			log.WithError(err).Error("error fetching API key")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		if !apiKeyObj.IsValid() {
			log.WithFields(log.Fields{
				"ip_address": ipAddress,
				"owner":      apiKeyObj.Owner,
			}).Warn("received request with invalid or expired API key")
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "API key is not valid or expired"})
			return
		}

		c.Set("admin_user", apiKeyObj.Owner)
		c.Next()
	}
}

// ClientIdMiddleware extracts the client ID from the X-Client-Id header
// and stores it in the request context for use by handlers.
func ClientIdMiddleware(db Persistence) gin.HandlerFunc {
	return func(c *gin.Context) {
		ipAddress := c.ClientIP()

		id := c.Request.Header.Get("X-Client-Id")
		if id == "" {
			log.WithFields(log.Fields{
				"ip_address": ipAddress,
			}).Warn("received request without client id")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Client ID is required"})
			return
		}

		client, err := db.GetClient(id)
		var errClientNotFound ErrClientNotFound
		if err != nil && errors.As(err, &errClientNotFound) {
			log.WithFields(log.Fields{
				"ip_address": ipAddress,
			}).Warn("received request with unknown client id")
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid client ID"})
			return
		} else if err != nil {
			log.WithError(err).Error("error fetching client")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.Set("client_id", client.Id)
		c.Next()
	}
}

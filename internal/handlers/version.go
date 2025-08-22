package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIVersion represents the current API version information
type APIVersion struct {
	Version     string `json:"version"`
	BuildTime   string `json:"build_time,omitempty"`
	GitCommit   string `json:"git_commit,omitempty"`
	Environment string `json:"environment"`
}

// VersionHandler handles version-related requests
type VersionHandler struct {
	version APIVersion
}

// NewVersionHandler creates a new version handler
func NewVersionHandler(version, buildTime, gitCommit, environment string) *VersionHandler {
	return &VersionHandler{
		version: APIVersion{
			Version:     version,
			BuildTime:   buildTime,
			GitCommit:   gitCommit,
			Environment: environment,
		},
	}
}

// GetVersion returns the current API version
// @Summary Get API version information
// @Description Returns the current API version, build time, and other metadata
// @Tags system
// @Produce json
// @Success 200 {object} APIVersion
// @Router /version [get]
func (h *VersionHandler) GetVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    h.version,
	})
}

// GetAPIInfo returns comprehensive API information
// @Summary Get comprehensive API information
// @Description Returns API version, supported versions, deprecation notices, etc.
// @Tags system
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/info [get]
func (h *VersionHandler) GetAPIInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"current_version": h.version,
			"supported_versions": []string{
				"v1",
			},
			"deprecated_versions": []string{},
			"api_documentation":   "/swagger/index.html",
			"endpoints": gin.H{
				"health":  "/health",
				"version": "/version",
				"v1":      "/api/v1",
			},
		},
	})
}

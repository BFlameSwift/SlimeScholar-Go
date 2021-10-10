package v1

import (


	"github.com/gin-gonic/gin"

	"net/http"
)

// Index doc
// @description 测试 Index 页
// @Tags 测试
// @Success 200 {string} string "{"success": true, "message": "gcp"}"
// @Router                                                                                                                                           /api/v1 [GET]
func Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "gcp"})
}

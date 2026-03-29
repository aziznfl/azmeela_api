package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParseID(c *gin.Context, param string) (int, bool) {
	idStr := c.Param(param)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return 0, false
	}
	return id, true
}

func GetPaginationParams(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	return page, limit
}

func BindJSON(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payload request tidak valid"})
		return false
	}
	return true
}

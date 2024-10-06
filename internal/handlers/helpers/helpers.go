package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandle(code int, err error) (c *gin.Context) {
	if err != nil {
		//логируем
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}
	return
}

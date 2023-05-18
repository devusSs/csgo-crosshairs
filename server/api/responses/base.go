package responses

import "github.com/gin-gonic/gin"

type SuccessResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

type ErrorResponse struct {
	Code  int `json:"code"`
	Error struct {
		ErrorCode    string `json:"error_code"`
		ErrorMessage string `json:"error_message"`
	} `json:"error"`
}

func (resp *SuccessResponse) SendSuccessReponse(c *gin.Context) {
	c.JSON(resp.Code, resp)
}

func (resp *ErrorResponse) SendErrorResponse(c *gin.Context) {
	c.JSON(resp.Code, resp)
}

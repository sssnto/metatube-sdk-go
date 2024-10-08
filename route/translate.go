package route

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"

	"github.com/metatube-community/metatube-sdk-go/translate"
	_ "github.com/metatube-community/metatube-sdk-go/translate/baidu"
	_ "github.com/metatube-community/metatube-sdk-go/translate/deepl"
	_ "github.com/metatube-community/metatube-sdk-go/translate/google"
	_ "github.com/metatube-community/metatube-sdk-go/translate/googlefree"
	_ "github.com/metatube-community/metatube-sdk-go/translate/openai"
)

type translateQuery struct {
	Q      string `form:"q" binding:"required"`
	From   string `form:"from"`
	To     string `form:"to" binding:"required"`
	Engine string `form:"engine" binding:"required"`
}

func getTranslate() gin.HandlerFunc {
	decoder := schema.NewDecoder()
	decoder.SetAliasTag("json")
	decoder.IgnoreUnknownKeys(true)

	return func(c *gin.Context) {
		query := &translateQuery{
			From: "auto",
		}
		if err := c.ShouldBindQuery(query); err != nil {
			abortWithStatusMessage(c, http.StatusBadRequest, err)
			return
		}
		engine := strings.ToLower(query.Engine)

		config, err := translate.BuildConfig(engine, func(config any) error {
			return decoder.Decode(config, c.Request.URL.Query())
		})
		if err != nil {
			abortWithStatusMessage(c, http.StatusBadRequest, err)
			return
		}

		result, err := translate.Translate(engine, query.Q, query.From, query.To, config)
		if err != nil {
			abortWithError(c, err)
			return
		}

		c.JSON(http.StatusOK, &responseMessage{
			Data: gin.H{
				"from":            query.From,
				"to":              query.To,
				"translated_text": result,
			},
		})
	}
}

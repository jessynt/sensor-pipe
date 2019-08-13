package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	sensorsanalytics "github.com/sensorsdata/sa-sdk-go"
	"github.com/tidwall/gjson"
)

func MakeTrackSignupHandler(sa sensorsanalytics.SensorsAnalytics) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rawData, err := ctx.GetRawData()
		if err != nil {
			panic(err)
		}

		if !strings.HasPrefix(ctx.GetHeader("Content-Type"), "application/json") {
			ctx.AbortWithStatus(400)
			return
		}

		rawString := string(rawData)

		rDistinctId := gjson.Get(rawString, "distinct_id")
		if !rDistinctId.Exists() {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
				"err_code": 61000002,
				"msg":      "is_login_id is required",
			})
			return
		}

		rOriginId := gjson.Get(rawString, "origin_id")
		if !rOriginId.Exists() {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
				"err_code": 61000002,
				"msg":      "origin_id is required",
			})
			return
		}
		err = sa.TrackSignup(rDistinctId.String(), rOriginId.String())
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
				"err_code": 61000002,
				"msg":      err.Error(),
			})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
			"err_code": 0,
			"msg":      "",
		})
	}
}

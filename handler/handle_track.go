package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	sensorsanalytics "github.com/sensorsdata/sa-sdk-go"
	"github.com/tidwall/gjson"
)

func MakeTrackHandler(sa sensorsanalytics.SensorsAnalytics) gin.HandlerFunc {
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
				"msg":      "distinct_id is required",
			})
			return
		}

		rType := gjson.Get(rawString, "event_type")
		if !rType.Exists() {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
				"err_code": 61000002,
				"msg":      "event_type is required",
			})
			return
		}

		rEvent := gjson.Get(rawString, "event_name")
		if rType.String() == "track" && !rEvent.Exists() {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
				"err_code": 61000002,
				"msg":      "event_name is required",
			})
			return
		}

		rIsLoginId := gjson.Get(rawString, "is_login_id")
		if !rIsLoginId.Exists() {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
				"err_code": 61000002,
				"msg":      "is_login_id is required",
			})
			return
		}

		rProperties := gjson.Get(rawString, "properties")
		if !rProperties.Exists() {
			ctx.AbortWithStatus(400)
			return
		}

		properties := make(map[string]interface{})

		for k, v := range rProperties.Map() {
			switch v.Type {
			case gjson.True:
				properties[k] = true
			case gjson.False:
				properties[k] = false
			case gjson.Number:
				properties[k] = v.Float()
			case gjson.String:
				properties[k] = v.String()
			default:
				ctx.AbortWithStatus(400)
				return
			}
		}

		switch rType.String() {
		case sensorsanalytics.TRACK:
			err = sa.Track(rDistinctId.String(), rEvent.String(), properties, rIsLoginId.Bool())
		case sensorsanalytics.PROFILE_SET:
			err = sa.ProfileSet(rDistinctId.String(), properties, rIsLoginId.Bool())
		case sensorsanalytics.PROFILE_SET_ONCE:
			err = sa.ProfileSetOnce(rDistinctId.String(), properties, rIsLoginId.Bool())
		default:
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
				"err_code": 61000002,
				"msg":      "event type not exists",
			})
			return
		}

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
				"err_code": 61000001,
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

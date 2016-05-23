package handler

import (
	"encoding/base64"
	"fmt"
	"github.com/Comcast/webpa-common/convey"
	"github.com/Comcast/webpa-common/fact"
	"golang.org/x/net/context"
	"net/http"
)

func Convey() ChainHandler {
	return ConveyCustom(ConveyHeader, base64.StdEncoding)
}

func ConveyCustom(conveyHeader string, encoding *base64.Encoding) ChainHandler {
	return ChainHandlerFunc(func(ctx context.Context, response http.ResponseWriter, request *http.Request, next ContextHandler) {
		rawPayload := request.Header.Get(conveyHeader)
		if len(rawPayload) > 0 {
			// BUG: https://www.teamccp.com/jira/browse/WEBPA-787
			const notAvailable string = "not-available"
			if rawPayload == notAvailable {
				fact.MustLogger(ctx).Error("Invalid convey header: %s.  FIX ME: https://www.teamccp.com/jira/browse/WEBPA-787", rawPayload)
			} else if conveyPayload, err := convey.ParsePayload(encoding, rawPayload); err != nil {
				WriteJsonError(
					response,
					http.StatusBadRequest,
					fmt.Sprintf(InvalidConveyPattern, rawPayload, err),
				)

				return
			} else {
				ctx = fact.SetConvey(ctx, conveyPayload)
			}
		}

		next.ServeHTTP(ctx, response, request)
	})
}
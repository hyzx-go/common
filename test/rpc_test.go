package test

import (
	"context"
	"fmt"
	"github.com/hyzx-go/common-b2c/config"
	r "github.com/hyzx-go/common-b2c/rpc"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	ctx := context.Background()
	reqUrl := "http://manage.oaid.com.cn/admin/user/page"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOjEsInN0YXR1cyI6MiwidXNlcm5hbWUiOiJhZG1pbiIsInBhc3N3b3JkIjoiOGRkY2ZmM2E4MGY0MTg5Y2ExYzlkNGQ5MDJjM2M5MDkiLCJlbWFpbCI6Ijg4ODg4ODg4QHFxLmNvbSIsInBob25lIjoiMTg4ODg4ODg4IiwiY3JlYXRlZEF0IjoiMjAyNC0xMS0wN1QxNjo1NzoxMSswODowMCIsInVwZGF0ZWRBdCI6IjIwMjQtMTEtMjdUMTA6NDU6MTArMDg6MDAiLCJpc3MiOiJoaSB0ZWNoIiwic3ViIjoibWFuYWdlIGJhY2tzdGFnZSB1c2VyIiwiZXhwIjoxNzMyNzc2MjE1LCJuYmYiOjE3MzI2ODk4MTUsImlhdCI6MTczMjY4OTgxNSwianRpIjoiMSJ9.2c-o_2UAg5vqFW_GEEbej649PJyChNxwGsGF38lo2j4"
	headers := r.Headers{r.Authorization: token}
	reqDTO := r.NewHttpClientBuilder().SetRequestType(r.Get).SetUrl(reqUrl).SetHeaders(headers).SetPrintLog(true).Build()

	data, err := config.GetParser().GetHTTPClient().Sync(ctx, reqDTO, time.Minute*3)
	if err != nil {
		t.Log("Sync err:", err.Error())
	}
	fmt.Println(string(data))
}

package complexpkg

import (
	"crypto/elliptic"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"time"
)

func useComplexPackages() {
	// see issue 25: https://github.com/nishanths/exhaustive/issues/25
	var (
		_ http.Server
		_ tls.Conn
		_ reflect.ChanDir
		_ json.Encoder
		_ elliptic.Curve
		_ time.Ticker
		_ os.File
	)
	fmt.Println(os.LookupEnv(""))
}

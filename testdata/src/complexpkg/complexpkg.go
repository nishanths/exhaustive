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
		_ elliptic.Curve
		_ tls.Conn
		_ json.Encoder
		_ fmt.Formatter
		_ http.Server
		_ os.File
		_ reflect.ChanDir
		_ time.Ticker
	)
	fmt.Println(os.Getgid())
}

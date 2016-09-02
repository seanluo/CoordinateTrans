# CoordinateTrans
Transform coordinates from/to China's offset version to WGS.

Two way transformation: 

	WGS <==> GCJ
	WGS <==> Baidu Maps
Usage: `go get github.com/seanluo/CoordinateTrans/trans`

hello_forbidden_city.go:
```go
package main

import (
	"fmt"
	"github.com/seanluo/CoordinateTrans/trans"
)

func main() {
	offset_location := trans.Location{
		Lat: 39.920082,
		Lng: 116.403609,
	}
	un_offset := trans.Bd2wgs(offset_location)
	fmt.Println(un_offset)
}
```


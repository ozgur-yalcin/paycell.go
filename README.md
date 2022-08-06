# Paycell.go
Turkcell (Paycell) API with golang

# Installation
```bash
go get github.com/ozgur-soft/paycell.go
```

# Satış işlemi
```go
package main

import (
	"context"
	"encoding/json"
	"fmt"

	paycell "github.com/ozgur-soft/paycell.go/src"
)

// Pos bilgileri
const (
	envmode  = "TEST"                // Çalışma ortamı (Production : "PROD" - Test : "TEST")
	merchant = "9998"                // İşyeri numarası
	storekey = "PAYCELL12345"        // İşyeri anahtarı
	appname  = "PAYCELLTEST"         // Uygulama adı
	apppass  = "PaycellTestPassword" // Uygulama şifresi
	prefix   = "666"                 // Referans no ilk 3 hanesi
)

func main() {
	api, req := paycell.Api(merchant, apppass, appname)
	api.SetStoreKey(storekey)
	api.SetPrefix(prefix)
	api.SetMode(envmode)
	api.SetIPAddress("127.0.0.1")         // IP adresi (zorunlu)
	api.SetPhoneNumber("905591111177")    // Müşteri numarası (zorunlu)
	api.SetAmount("1.00", "TRY")          // Satış tutarı (zorunlu)
	req.SetCardNumber("4355084355084358") // Kart numarası (zorunlu)
	req.SetCardExpiry("12", "26")         // Son kullanma tarihi - AA,YY (zorunlu)
	req.SetCardCode("000")                // Kart arkasındaki 3 haneli numara (zorunlu)
	req.Provision.Installment = "0"       // Taksit sayısı (varsa)

	ctx := context.Background()
	if res, err := api.Auth(ctx, req); err == nil {
		pretty, _ := json.MarshalIndent(res.Provision, " ", " ")
		fmt.Println(string(pretty))
	} else {
		fmt.Println(err)
	}
}
```

# İade işlemi
```go
package main

import (
	"context"
	"encoding/json"
	"fmt"

	paycell "github.com/ozgur-soft/paycell.go/src"
)

// Pos bilgileri
const (
	envmode  = "TEST"                // Çalışma ortamı (Production : "PROD" - Test : "TEST")
	merchant = "9998"                // İşyeri numarası
	storekey = "PAYCELL12345"        // İşyeri anahtarı
	appname  = "PAYCELLTEST"         // Uygulama adı
	apppass  = "PaycellTestPassword" // Uygulama şifresi
	prefix   = "666"                 // Referans no ilk 3 hanesi
)

func main() {
	api, req := paycell.Api(merchant, apppass, appname)
	api.SetStoreKey(storekey)
	api.SetPrefix(prefix)
	api.SetMode(envmode)
	api.SetPhoneNumber("905591111177") // Müşteri numarası (zorunlu)
	api.SetIPAddress("127.0.0.1")      // IP adresi (zorunlu)
	api.SetAmount("1.00", "TRY")       // İade tutarı (zorunlu)

	req.Refund.OriginalRefNo = "" // Referans numarası (zorunlu)

	ctx := context.Background()
	if res, err := api.Refund(ctx, req); err == nil {
		pretty, _ := json.MarshalIndent(res.Refund, " ", " ")
		fmt.Println(string(pretty))
	} else {
		fmt.Println(err)
	}
}
```

# İptal işlemi
```go
package main

import (
	"context"
	"encoding/json"
	"fmt"

	paycell "github.com/ozgur-soft/paycell.go/src"
)

// Pos bilgileri
const (
	envmode  = "TEST"                // Çalışma ortamı (Production : "PROD" - Test : "TEST")
	merchant = "9998"                // İşyeri numarası
	storekey = "PAYCELL12345"        // İşyeri anahtarı
	appname  = "PAYCELLTEST"         // Uygulama adı
	apppass  = "PaycellTestPassword" // Uygulama şifresi
	prefix   = "666"                 // Referans no ilk 3 hanesi
)

func main() {
	api, req := paycell.Api(merchant, apppass, appname)
	api.SetStoreKey(storekey)
	api.SetPrefix(prefix)
	api.SetMode(envmode)
	api.SetPhoneNumber("905591111177") // Müşteri numarası (zorunlu)
	api.SetIPAddress("127.0.0.1")      // IP adresi (zorunlu)

	req.Cancel.OriginalRefNo = "" // Referans numarası (zorunlu)

	ctx := context.Background()
	if res, err := api.Cancel(ctx, req); err == nil {
		pretty, _ := json.MarshalIndent(res.Cancel, " ", " ")
		fmt.Println(string(pretty))
	} else {
		fmt.Println(err)
	}
}
```

# Mobil ödeme açma işlemi
```go
package main

import (
	"context"
	"encoding/json"
	"fmt"

	paycell "github.com/ozgur-soft/paycell.go/src"
)

// Pos bilgileri
const (
	envmode  = "TEST"                // Çalışma ortamı (Production : "PROD" - Test : "TEST")
	merchant = "9998"                // İşyeri numarası
	storekey = "PAYCELL12345"        // İşyeri anahtarı
	appname  = "PAYCELLTEST"         // Uygulama adı
	apppass  = "PaycellTestPassword" // Uygulama şifresi
	prefix   = "666"                 // Referans no ilk 3 hanesi
)

func main() {
	api, req := paycell.Api(merchant, apppass, appname)
	api.SetStoreKey(storekey)
	api.SetPrefix(prefix)
	api.SetMode(envmode)
	api.SetPhoneNumber("905305289290") // Müşteri numarası (zorunlu)
	api.SetIPAddress("127.0.0.1")      // IP adresi (zorunlu)

	ctx := context.Background()
	if get, err := api.GetPaymentMethods(ctx, req); err == nil {
		if get.PaymentMethods.MobilePayment != nil {
			if !get.PaymentMethods.MobilePayment.IsDcbOpen {
				switch get.PaymentMethods.MobilePayment.IsEulaExpired {
				case true: // Sözleşmesi Güncel Olmayan Müşteri İçin
					req.MobilePayment.EulaID = get.PaymentMethods.MobilePayment.EulaId
					if open, err := api.OpenMobilePayment(ctx, req); err == nil {
						pretty, _ := json.MarshalIndent(open.MobilePayment, " ", " ")
						fmt.Println(string(pretty))
					} else {
						fmt.Println(err)
					}
				case false: // Sözleşmesi Güncel Olan Müşteri İçin
					if open, err := api.OpenMobilePayment(ctx, req); err == nil {
						pretty, _ := json.MarshalIndent(open.MobilePayment, " ", " ")
						fmt.Println(string(pretty))
					} else {
						fmt.Println(err)
					}
				}
			}
		}
	} else {
		fmt.Println(err)
	}
}
```
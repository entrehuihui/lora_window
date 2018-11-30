package updata

import (
	"github.com/brocaar/lora-app-server/internal/config"
	"github.com/brocaar/lora-app-server/internal/handler"
	"github.com/brocaar/lorawan"
)

//DownData ...
func DownData() {
	dataDownPayload := handler.DataDownPayload{
		ApplicationID: 3,
		DevEUI:        lorawan.EUI64{36, 197, 217, 230, 50, 87, 245, 140},
		Confirmed:     true,
		FPort:         5453,
		Data:          []byte("I AM DOWN NOW"),
	}
	config.C.ApplicationServer.Integration.Handler.DataDownChan() <- dataDownPayload
}

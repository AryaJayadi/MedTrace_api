package web

import (
	"fmt"
	"net/http"
)

func (setup OrgSetup) Query(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Recieved Query request")
	queryparams := r.URL.Query()
	chainCodeName := queryparams.Get("chaincodeid")
	channelId := queryparams.Get("channelid")
	function := queryparams.Get("function")
	args := r.URL.Query()["args"]
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %v\n", channelId, chainCodeName, function, args)
	network := setup.Gateway.GetNetwork(channelId)
	contract := network.GetContract(chainCodeName)
	evaluateResponse, err := contract.EvaluateTransaction(function, args...)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		return
	}
	fmt.Fprintf(w, "Response: %s", evaluateResponse)
}

package utils

import (
	"fmt"
	"github.com/ndphu/betmonitor/config"
	"github.com/ndphu/message-handler-lib/broker"
	"log"
	"regexp"
)

var InsideSpacesRegex = regexp.MustCompile(`[\s\p{Zs}]{2,}`)

func Reply(result string) error {
	conf := config.GetConfig()

	for _, target := range conf.Targets {
		for _, thread := range target.Threads {
			rpc, err := broker.NewRpcClient(target.WorkerId)
			if err != nil {
				log.Printf("React: Fail to create RPC client by error %v\n", err)
				continue
			}
			request := &broker.RpcRequest{
				Method: "sendText",
				Args:   []string{thread, WrapAsPreformatted(result)},
			}
			if err := rpc.Send(request); err != nil {
				log.Println("React: Fail to reply:", err.Error())
				continue
			}
		}
	}

	return nil
}

func WrapAsPreformatted(message string) string {
	return fmt.Sprintf("<pre raw_pre=\"{code}\" raw_post=\"{code}\">%s</pre>", message)
}

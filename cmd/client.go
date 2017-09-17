package main

import yyyp_client "github.com/shellus/yyyp/p2p-client"

func main(){
	yclient, err := yyyp_client.NewClient()
	if err != nil {
		panic(err)
	}
	defer yclient.Close()
	err = yclient.RunClient()
	if err != nil {
		panic(err)
	}
}
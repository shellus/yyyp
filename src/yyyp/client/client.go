package client

func Connect(){
	&message{
		action:registerTunnel,
		data:&registerTunnelMessage{tunnelName:tunnelName},
	}
}
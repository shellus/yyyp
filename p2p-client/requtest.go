package p2p_client


import "github.com/shellus/yyyp/pack"


func (t *P2PClient) WritePack(message interface{})(err error) {
	t.debug("send pack [%#v]", message)

	body, err := pack.Package(message)
	if err != nil {
		return
	}
	err = t.YConn.WriteMessage(body)
	return
}
func (t *P2PClient) RequestLink(name string) {
	t.WritePack(pack.PackLink{Name: name})
	t.debug("send link request : [%s]", name)
}

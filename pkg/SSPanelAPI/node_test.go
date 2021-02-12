package SSPanelAPI

import (
	"fmt"
	"testing"
)

func TestParseRawInfo(t *testing.T) {
	serverRawInfo := "www.google.com;31824;2;tls;ws;path=\\/videos|host=download.windowsupdate.com|relay"
	node := &NodeInfo{
		RawInfo: serverRawInfo,
	}
	err := node.parseRawInfo()
	if err != nil {
		t.Errorf("Failed")
	}
	fmt.Println(*node)
}

package config

import (
	"fmt"

	"github.com/micromdm/micromdm/vpp"
)

const VPPTokenTopic = "mdm.VPPTokenAdded"

type VPPToken struct {
	UDID   string `json:"udid"`
	SToken string `json:"sToken"`
}

// create a VPP client from token.
func (tok VPPToken) Client() (*vpp.Client, error) {
	fmt.Printf("vppToken: %s\n", tok.SToken)
	client, err := vpp.NewClient(tok.SToken)
	if err != nil {
		fmt.Printf("vppToken err: %v\n", err)
		return nil, err
	}
	fmt.Printf("vppToken OK: %v\n", client)
	return client, nil
}

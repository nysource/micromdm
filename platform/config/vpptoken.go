package config

import (
	"log"

	"github.com/micromdm/micromdm/vpp"
)

const VPPTokenTopic = "mdm.VPPTokenAdded"

type VPPToken struct {
	UDID   string `json:"udid"`
	SToken string `json:"sToken"`
}

// create a VPP client from token.
func (tok VPPToken) Client() (*vpp.Client, error) {
	log.Printf("vppToken: %s\n", tok.SToken)
	client, err := vpp.NewClient(tok.SToken)
	if err != nil {
		log.Printf("vppToken err: %s\n", err)
		return nil, err
	}
	log.Printf("vppToken OK: %s\n", client)
	return client, nil
}

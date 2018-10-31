package config

import (
	//"time"

	"github.com/micromdm/micromdm/vpp"
)

const VPPTokenTopic = "mdm.VPPTokenAdded"

type VPPToken struct {
	UDID    string `json:"udid"`
	SToken  string `json:"sToken"`
	OrgName string `json:"orgName"`
	Token   string `json:"token"`
	ExpDate string `json:"expDate"`
}

// create a VPP client from token.
func (tok VPPToken) Client() (*vpp.Client, error) {
	client, err := vpp.NewClient(tok.SToken)
	if err != nil {
		return nil, err
	}
	return client, nil
}

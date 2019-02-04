package vpp

import (
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
)

type VPPAppsList struct {
	VPPApps []VPPApp `json:"vpp-apps"`
}

type VPPApp struct {
	ID                string   `json:"id"`
	DeviceFamilies    []string `json:"deviceFamilies"`
	Name              string   `json:"name"`
	TotalLicenses     int      `json:"totalLicenses"`
	AssignedLicenses  int      `json:"assignedLicenses"`
	AvailableLicenses int      `json:"availableLicenses"`
	TokenUUID         string   `json:"tokenUUID"`
}

func (c *Client) GetVPPApps() (*VPPAppsList, error) {

	assetsSrvOptions := AssetsSrvOptions{
		SToken:               c.VPPToken.SToken,
		IncludeLicenseCounts: true,
	}

	assetsSrv, err := c.GetAssetsSrv(assetsSrvOptions)
	if err != nil {
		return nil, err
	}

	var clientContext ClientContext
	context := assetsSrv.ClientContext

	err = json.NewDecoder(strings.NewReader(context)).Decode(&clientContext)
	if err != nil {
		return nil, errors.Wrap(err, "decode ClientContext")
	}

	appsList := VPPAppsList{}

	var assets = assetsSrv.Assets
	for _, asset := range assets {

		contentMetadataOptions := ContentMetadataOptions{
			ID: asset.AdamIDStr,
		}
		data, err := c.GetAppData(contentMetadataOptions)
		if err != nil {
			return nil, err
		}

		vppApp := VPPApp{
			ID:                asset.AdamIDStr,
			DeviceFamilies:    data.DeviceFamilies,
			Name:              data.Name,
			TotalLicenses:     asset.TotalCount,
			AssignedLicenses:  asset.AssignedCount,
			AvailableLicenses: asset.AvailableCount,
			TokenUUID:         clientContext.GUID,
		}
		appsList.VPPApps = append(appsList.VPPApps, vppApp)
	}

	return &appsList, errors.Wrap(err, "make LicensesSrv request")
}

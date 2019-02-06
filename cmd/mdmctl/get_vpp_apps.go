package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/micromdm/micromdm/vpp"
)

type vppAppsTableOutput struct{ w *tabwriter.Writer }

func (out *vppAppsTableOutput) BasicHeader() {
	fmt.Fprintf(out.w, "ID\tName\tDeviceFamilies\tLicenses\tAssigned\tAvailable\tTokenUUID\n")
}

func (out *vppAppsTableOutput) BasicFooter() {
	out.w.Flush()
}

type VPPAppResponse struct {
	Asset         interface{} `json:"vpp-asset,omitempty"`
	ClientContext string      `json:"client-context,omitempty"`
	Metadata      interface{} `json:"vpp-metadata,omitempty"`
}

func (cmd *getCommand) getVPPApps(args []string) error {
	flagset := flag.NewFlagSet("vpp-apps", flag.ExitOnError)
	var (
		flIDFilter = flagset.String("id", "", "specify the id of the vpp app to get full details")
	)
	flagset.Usage = usageFor(flagset, "mdmctl get vpp-apps [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}
	ctx := context.Background()

	if *flIDFilter != "" {
		var response VPPAppResponse

		assetsSrv, err := cmd.vppsvc.GetAssetsSrv(ctx, vpp.AssetsSrvOptions{
			IncludeLicenseCounts: true,
		})
		if err != nil {
			return err
		}
		response.ClientContext = assetsSrv.ClientContext
		assets := assetsSrv.Assets

		for _, a := range assets {
			if a.AdamIDStr == *flIDFilter {
				response.Asset = a

				metadata, err := cmd.vppsvc.GetContentMetadata(ctx, vpp.ContentMetadataOptions{
					ID: *flIDFilter,
				})
				if err != nil {
					return err
				}
				response.Metadata = metadata
			}
		}

		bytes, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(bytes))
		return nil
	}

	appsList, err := cmd.vppsvc.GetVPPApps(ctx)
	apps := appsList.VPPApps
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	out := vppAppsTableOutput{w}
	out.BasicHeader()
	defer out.BasicFooter()
	for _, a := range apps {
		fmt.Fprintf(out.w, "%s\t%s\t%s\t%d\t%d\t%d\t%s\n", a.ID, a.Name, a.DeviceFamilies, a.TotalLicenses, a.AssignedLicenses, a.AvailableLicenses, a.TokenUUID)
	}
	return nil
}

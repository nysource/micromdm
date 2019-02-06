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
	VPPAppsList []VPPAppData `json:"vpp-apps"`
}

type VPPAppData struct {
	Asset         interface{}       `json:"vpp-asset,omitempty"`
	ClientContext vpp.ClientContext `json:"client-context,omitempty"`
	Metadata      interface{}       `json:"vpp-metadata,omitempty"`
}

func (cmd *getCommand) getVPPApps(args []string) error {
	flagset := flag.NewFlagSet("vpp-apps", flag.ExitOnError)
	var (
		flIDFilter = flagset.String("id", "", "specify the id of the vpp app to filter results to one app")
		flVerbose  = flagset.Bool("verbose", false, "specify -verbose to get full apps details")
	)
	flagset.Usage = usageFor(flagset, "mdmctl get vpp-apps [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}
	ctx := context.Background()

	if *flVerbose == true {
		assetsSrv, err := cmd.vppsvc.GetAssetsSrv(ctx, vpp.AssetsSrvOptions{
			IncludeLicenseCounts: true,
		})
		if err != nil {
			return err
		}
		assets := assetsSrv.Assets

		if *flIDFilter != "" {
			for _, a := range assets {
				if a.AdamIDStr == *flIDFilter {
					assets = []vpp.Asset{a}
				}
			}
			if len(assets) > 1 {
				assets = []vpp.Asset{}
			}
		}

		var response VPPAppResponse
		for _, a := range assets {
			var vppApp VPPAppData

			var clientContext vpp.ClientContext
			vpp.DecodeToClientContext(assetsSrv.ClientContext, &clientContext)

			vppApp.ClientContext = clientContext
			vppApp.Asset = a

			metadata, err := cmd.vppsvc.GetContentMetadata(ctx, vpp.ContentMetadataOptions{
				ID: a.AdamIDStr,
			})
			if err != nil {
				return err
			}
			vppApp.Metadata = metadata

			response.VPPAppsList = append(response.VPPAppsList, vppApp)
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

	if *flIDFilter != "" {
		for _, a := range apps {
			if a.ID == *flIDFilter {
				apps = []vpp.VPPApp{a}
			}
		}
		if len(apps) > 1 {
			apps = []vpp.VPPApp{}
		}
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

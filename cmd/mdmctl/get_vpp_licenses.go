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

type vppLicenseWithSerials struct {
	AdamID  string
	Serials []string
}

type vppLicensesTableOutput struct{ w *tabwriter.Writer }

func (out *vppLicensesTableOutput) BasicHeader() {
	fmt.Fprintf(out.w, "AppID\tTokenUUID\tSerials\n")
}

func (out *vppLicensesTableOutput) BasicFooter() {
	out.w.Flush()
}

func (cmd *getCommand) getVPPLicenses(args []string) error {
	flagset := flag.NewFlagSet("vpp-licenses", flag.ExitOnError)
	var (
		flIDFilter = flagset.String("id", "", "specify the id of the vpp app to get full details")
	)
	flagset.Usage = usageFor(flagset, "mdmctl get vpp-licenses [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}
	ctx := context.Background()

	var options vpp.LicensesSrvOptions
	if *flIDFilter != "" {
		options.AdamID = *flIDFilter
	}

	licensesSrv, err := cmd.vppsvc.GetLicensesSrv(ctx, options)
	if err != nil {
		return err
	}

	if *flIDFilter != "" {
		bytes, err := json.MarshalIndent(licensesSrv, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(bytes))
		return nil
	}

	/*var clientContext vpp.ClientContext
	context := licensesSrv.ClientContext

	err = json.NewDecoder(strings.NewReader(context)).Decode(&clientContext)
	if err != nil {
		return err
	}*/

	licenses := licensesSrv.Licenses

	var licensesList []vppLicenseWithSerials
	for _, license := range licenses {
		found := false
		if len(licensesList) > 0 {
			for _, saved := range licensesList {
				if license.AdamIDStr == saved.AdamID {
					found = true
					saved.Serials = append(saved.Serials, license.SerialNumber)
				}
			}
		}

		if !found {
			licensesList = append(licensesList, vppLicenseWithSerials{
				AdamID:  license.AdamIDStr,
				Serials: []string{license.SerialNumber},
			})
		}
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	out := vppLicensesTableOutput{w}
	out.BasicHeader()
	defer out.BasicFooter()
	for _, a := range licensesList {
		//fmt.Fprintf(out.w, "%s\t%s\t%s\n", a.AdamID, clientContext.GUID, a.Serials)
		fmt.Fprintf(out.w, "%s\t%s\t%s\n", a.AdamID, "", a.Serials)
	}
	return nil
}

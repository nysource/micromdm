package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
)

type vppAppsTableOutput struct{ w *tabwriter.Writer }

func (out *vppAppsTableOutput) BasicHeader() {
	fmt.Fprintf(out.w, "ID\tName\tDeviceFamilies\tLicenses\tAssigned\tAvailable\tTokenUUID\n")
}

func (out *vppAppsTableOutput) BasicFooter() {
	out.w.Flush()
}

func (cmd *getCommand) getVPPApps(args []string) error {
	flagset := flag.NewFlagSet("vpp-apps", flag.ExitOnError)
	var (
		flIDFilter    = flagset.String("id", "", "specify the id of the vpp app to get full details")
		flTokenFilter = flagset.String("token", "", "specify the token uuid to filter by token")
	)
	flagset.Usage = usageFor(flagset, "mdmctl get vpp-apps [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}
	ctx := context.Background()
	appsList, err := cmd.vppsvc.GetVPPApps(ctx)
	apps := appsList.VPPApps

	/*apps, err := cmd.vppsvc.GetVPPApps(ctx, appstore.ListAppsOption{
		FilterName: []string{*flIDFilter},
	})*/

	if err != nil {
		return err
	}

	if *flIDFilter != "" && (len(apps) != 0) {
		//payload := apps[0]
		//fmt.Println(string(payload))
		return nil
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

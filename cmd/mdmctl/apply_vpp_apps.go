package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"

	"github.com/micromdm/micromdm/vpp"
	"github.com/pkg/errors"
)

func (cmd *applyCommand) applyVPPApp(args []string) error {
	flagset := flag.NewFlagSet("vpp-apps", flag.ExitOnError)
	var (
		flAppID   = flagset.String("id", "", "specify the id of the vpp app for which to associate serials")
		flSerials = flagset.String("serials", "", "specify a list of serials to associate to the vpp app")
	)
	flagset.Usage = usageFor(flagset, "mdmctl apply vpp-apps [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}
	ctx := context.Background()

	if *flAppID == "" || *flSerials == "" {
		flagset.Usage()
		return errors.New("bad input: must provide an app id and at least one serial")
	}

	manage, err := cmd.vppsvc.ManageVPPLicensesByAdamIdSrv(ctx, vpp.ManageVPPLicensesByAdamIdSrvOptions{
		AdamIDStr:              *flAppID,
		AssociateSerialNumbers: []string{*flSerials},
	})
	if err != nil {
		return err
	}

	bytes, err := json.MarshalIndent(manage, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(bytes))
	return nil
}

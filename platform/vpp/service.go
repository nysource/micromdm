package vpp

import (
	"context"
	"sync"

	"github.com/micromdm/micromdm/platform/pubsub"
	"github.com/micromdm/micromdm/vpp"
)

type Service interface {
	GetVPPApps(ctx context.Context, ids []string) (*vpp.VPPAppResponse, error)
	GetContentMetadata(ctx context.Context, options vpp.ContentMetadataOptions) (*vpp.ContentMetadata, error)
	GetVPPAssetsSrv(ctx context.Context) (*vpp.VPPAssetsSrv, error)
	GetLicensesSrv(ctx context.Context, options vpp.GetLicensesSrvOptions) (*vpp.LicensesSrv, error)
	GetVPPServiceConfigSrv(ctx context.Context) (*vpp.VPPServiceConfigSrv, error)
	ManageVPPLicensesByAdamIdSrv(ctx context.Context, options vpp.ManageVPPLicensesByAdamIdSrvOptions) (*vpp.ManageVPPLicensesByAdamIdSrv, error)
}

type VPPClient interface {
	GetVPPApps(...string) (*vpp.VPPAppResponse, error)
	GetContentMetadata(vpp.ContentMetadataOptions) (*vpp.ContentMetadata, error)
	GetVPPAssetsSrv() (*vpp.VPPAssetsSrv, error)
	GetLicensesSrv(vpp.GetLicensesSrvOptions) (*vpp.LicensesSrv, error)
	GetVPPServiceConfigSrv() (*vpp.VPPServiceConfigSrv, error)
	ManageVPPLicensesByAdamIdSrv(string, vpp.ManageVPPLicensesByAdamIdSrvOptions) (*vpp.ManageVPPLicensesByAdamIdSrv, error)
}

type VPPService struct {
	mtx        sync.RWMutex
	client     VPPClient
	subscriber pubsub.Subscriber
}

func (svc *VPPService) Run() error {
	return svc.watchTokenUpdates(svc.subscriber)
}

func New(client VPPClient, subscriber pubsub.Subscriber) *VPPService {
	return &VPPService{client: client, subscriber: subscriber}
}

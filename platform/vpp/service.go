package vpp

import (
	"context"
	"sync"

	"github.com/micromdm/micromdm/vpp"
	"github.com/micromdm/micromdm/platform/pubsub"
)

type Service interface {
	GetLicensesSrv(ctx context.Context, options vpp.LicensesSrvOptions) (*vpp.LicensesSrv, error)
	GetVPPServiceConfigSrv(ctx context.Context, options vpp.VPPServiceConfigSrvOptions) (*vpp.VPPServiceConfigSrv, error)
}

type VPPClient interface {
	GetLicensesSrv(vpp.LicensesSrvOptions) (*vpp.LicensesSrv, error)
	GetVPPServiceConfigSrv(vpp.VPPServiceConfigSrvOptions) (*vpp.VPPServiceConfigSrv, error)
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

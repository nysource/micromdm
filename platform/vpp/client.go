package vpp

import (
	"net/url"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/micromdm/micromdm/pkg/httputil"
)

func NewHTTPClient(instance, token string, logger log.Logger, opts ...httptransport.ClientOption) (Service, error) {
	u, err := url.Parse(instance)
	if err != nil {
		return nil, err
	}

	var getLicensesSrvEndpoint endpoint.Endpoint
	{
		getLicensesSrvEndpoint = httptransport.NewClient(
			"POST",
			httputil.CopyURL(u, "/v1/vpp/licensessrv"),
			httputil.EncodeRequestWithToken(token, httptransport.EncodeJSONRequest),
			decodeGetLicensesSrvResponse,
			opts...,
		).Endpoint()
	}

	var getVPPServiceConfigSrvEndpoint endpoint.Endpoint
	{
		getVPPServiceConfigSrvEndpoint = httptransport.NewClient(
			"POST",
			httputil.CopyURL(u, "/v1/vpp/serviceconfigsrv"),
			httputil.EncodeRequestWithToken(token, httptransport.EncodeJSONRequest),
			decodeGetVPPServiceConfigSrvResponse,
			opts...,
		).Endpoint()
	}

	return Endpoints{
		GetLicensesSrvEndpoint:         getLicensesSrvEndpoint,
		GetVPPServiceConfigSrvEndpoint: getVPPServiceConfigSrvEndpoint,
	}, nil
}

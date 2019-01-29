package vpp

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-kit/kit/endpoint"

	"github.com/micromdm/micromdm/pkg/httputil"
	"github.com/micromdm/micromdm/vpp"
)

func (svc *VPPService) GetVPPServiceConfigSrv(ctx context.Context, options vpp.VPPServiceConfigSrvOptions) (*vpp.VPPServiceConfigSrv, error) {
	if svc.client == nil {
		return nil, errors.New("VPP not configured yet. add a VPP token to enable VPP")
	}
	return svc.client.GetVPPServiceConfigSrv(options)
}

//type getVPPServiceConfigSrvRequest struct{}
type getVPPServiceConfigSrvResponse struct {
	*vpp.VPPServiceConfigSrv
	Err error `json:"err,omitempty"`
}

func (r getVPPServiceConfigSrvResponse) Failed() error { return r.Err }

func decodeGetVPPServiceConfigSrvRequest(ctx context.Context, r *http.Request) (interface{}, error) {

	var req vpp.VPPServiceConfigSrvOptions

	switch r.Method {
	case "POST":
		err := httputil.DecodeJSONRequest(r, &req)
		return req, err
	default:
		return req, nil
	}
}

func decodeGetVPPServiceConfigSrvResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp getVPPServiceConfigSrvResponse
	err := httputil.DecodeJSONResponse(r, &resp)
	return resp, err
}

func MakeGetVPPServiceConfigSrvEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(vpp.VPPServiceConfigSrvOptions)
		configsrv, err := svc.GetVPPServiceConfigSrv(ctx, req)
		return getVPPServiceConfigSrvResponse{VPPServiceConfigSrv: configsrv, Err: err}, nil
	}
}

func (e Endpoints) GetVPPServiceConfigSrv(ctx context.Context, options vpp.VPPServiceConfigSrvOptions) (*vpp.VPPServiceConfigSrv, error) {
	response, err := e.GetVPPServiceConfigSrvEndpoint(ctx, options)
	if err != nil {
		return nil, err
	}
	return response.(getVPPServiceConfigSrvResponse).VPPServiceConfigSrv, response.(getVPPServiceConfigSrvResponse).Err
}

package vpp

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-kit/kit/endpoint"

	"github.com/micromdm/micromdm/pkg/httputil"
	"github.com/micromdm/micromdm/vpp"
)

func (svc *VPPService) GetVPPServiceConfigSrv(ctx context.Context) (*vpp.VPPServiceConfigSrv, error) {
	if svc.client == nil {
		return nil, errors.New("VPP not configured yet. add a VPP token to enable VPP")
	}
	return svc.client.GetVPPServiceConfigSrv()
}

type getVPPServiceConfigSrvRequest struct {
	SToken string `json:"sToken"`
}

type getServiceConfigSrvResponse struct {
	*vpp.VPPServiceConfigSrv
	Err error `json:"err,omitempty"`
}

func (r getServiceConfigSrvResponse) Failed() error { return r.Err }

func decodeGetServiceConfigSrvRequest(ctx context.Context, r *http.Request) (interface{}, error) {

	var req getVPPServiceConfigSrvRequest

	switch r.Method {
	case "POST":
		err := httputil.DecodeJSONRequest(r, &req)
		return req, err
	default:
		return req, nil
	}
}

func decodeGetServiceConfigSrvResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp getServiceConfigSrvResponse
	err := httputil.DecodeJSONResponse(r, &resp)
	return resp, err
}

func MakeGetVPPServiceConfigSrvEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		configsrv, err := svc.GetVPPServiceConfigSrv(ctx)
		return getServiceConfigSrvResponse{VPPServiceConfigSrv: configsrv, Err: err}, nil
	}
}

func (e Endpoints) GetVPPServiceConfigSrv(ctx context.Context) (*vpp.VPPServiceConfigSrv, error) {
	response, err := e.GetVPPServiceConfigSrvEndpoint(ctx, nil)
	if err != nil {
		return nil, err
	}
	return response.(getServiceConfigSrvResponse).VPPServiceConfigSrv, response.(getServiceConfigSrvResponse).Err
}

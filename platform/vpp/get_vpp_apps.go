package vpp

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-kit/kit/endpoint"

	"github.com/micromdm/micromdm/pkg/httputil"
	"github.com/micromdm/micromdm/vpp"
)

func (svc *VPPService) GetVPPApps(ctx context.Context) (*vpp.VPPAppsList, error) {
	if svc.client == nil {
		return nil, errors.New("VPP not configured yet. add a VPP token to enable VPP")
	}
	return svc.client.GetVPPApps()
}

type getVPPAppsResponse struct {
	*vpp.VPPAppsList
	Err error `json:"err,omitempty"`
}

func (r getVPPAppsResponse) Failed() error { return r.Err }

func decodeGetVPPAppsRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeGetVPPAppsResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp getVPPAppsResponse
	err := httputil.DecodeJSONResponse(r, &resp)
	return resp, err
}

func MakeGetVPPAppsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		srv, err := svc.GetVPPApps(ctx)
		return getVPPAppsResponse{VPPAppsList: srv, Err: err}, nil
	}
}

func (e Endpoints) GetVPPApps(ctx context.Context) (*vpp.VPPAppsList, error) {
	var options interface{}
	response, err := e.GetVPPAppsEndpoint(ctx, options)
	if err != nil {
		return nil, err
	}
	return response.(getVPPAppsResponse).VPPAppsList, response.(getVPPAppsResponse).Err
}

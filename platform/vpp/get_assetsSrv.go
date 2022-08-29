package vpp

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-kit/kit/endpoint"

	"github.com/micromdm/micromdm/pkg/httputil"
	"github.com/micromdm/micromdm/vpp"
)

func (svc *VPPService) GetVPPAssetsSrv(ctx context.Context) (*vpp.VPPAssetsSrv, error) {
	if svc.client == nil {
		return nil, errors.New("VPP not configured yet. add a VPP token to enable VPP")
	}
	return svc.client.GetVPPAssetsSrv()
}

type getVPPAssetsSrvRequest struct {
	SToken string `json:"sToken"`
}

type getVPPAssetsSrvResponse struct {
	*vpp.VPPAssetsSrv
	Err error `json:"err,omitempty"`
}

func (r getVPPAssetsSrvResponse) Failed() error { return r.Err }

func decodeGetVPPAssetsSrvRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req getVPPAssetsSrvRequest

	switch r.Method {
	case "POST":
		err := httputil.DecodeJSONRequest(r, &req)
		return req, err
	default:
		return req, nil
	}
}

func decodeGetVPPAssetsSrvResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp getVPPAssetsSrvResponse
	err := httputil.DecodeJSONResponse(r, &resp)
	return resp, err
}

func MakeGetVPPAssetsSrvEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		srv, err := svc.GetVPPAssetsSrv(ctx)
		return getVPPAssetsSrvResponse{VPPAssetsSrv: srv, Err: err}, nil
	}
}

func (e Endpoints) GetVPPAssetsSrv(ctx context.Context) (*vpp.VPPAssetsSrv, error) {
	response, err := e.GetVPPAssetsSrvEndpoint(ctx, nil)
	if err != nil {
		return nil, err
	}
	return response.(getVPPAssetsSrvResponse).VPPAssetsSrv, response.(getVPPAssetsSrvResponse).Err
}

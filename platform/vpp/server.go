package vpp

import (
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/micromdm/micromdm/pkg/httputil"
)

type Endpoints struct {
	GetLicensesSrvEndpoint    		endpoint.Endpoint
	GetVPPServiceConfigSrvEndpoint   endpoint.Endpoint
}

func MakeServerEndpoints(s Service, outer endpoint.Middleware, others ...endpoint.Middleware) Endpoints {
	return Endpoints{
		GetLicensesSrvEndpoint:         endpoint.Chain(outer, others...)(MakeGetLicensesSrvEndpoint(s)),
		GetVPPServiceConfigSrvEndpoint:   endpoint.Chain(outer, others...)(MakeGetVPPServiceConfigSrvEndpoint(s)),
	}
}

func RegisterHTTPHandlers(r *mux.Router, e Endpoints, options ...httptransport.ServerOption) {
	// POST		/v1/vpp/apps		              list all vpp apps
	// PUT		/v1/vpp/apps			            assign or unassign device licenses
	// POST		/v1/vpp/licensessrv		        list all vpp licenses
	// GET		/v1/vpp/vppserviceconfigsrv		get vppserviceconfigsrv information
	// POST		/v1/vpp/vppserviceconfigsrv		get vppserviceconfigsrv information for specific sToken

	r.Methods("POST").Path("/v1/vpp/licensessrv").Handler(httptransport.NewServer(
		e.GetLicensesSrvEndpoint,
		decodeGetLicensesSrvRequest,
		httputil.EncodeJSONResponse,
		options...,
	))

	r.Methods("GET").Path("/v1/vpp/vppserviceconfigsrv").Handler(httptransport.NewServer(
		e.GetVPPServiceConfigSrvEndpoint,
		decodeGetVPPServiceConfigSrvRequest,
		httputil.EncodeJSONResponse,
		options...,
	))

	r.Methods("POST").Path("/v1/vpp/vppserviceconfigsrv").Handler(httptransport.NewServer(
		e.GetVPPServiceConfigSrvEndpoint,
		decodeGetVPPServiceConfigSrvRequest,
		httputil.EncodeJSONResponse,
		options...,
	))
}

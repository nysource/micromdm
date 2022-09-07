package config

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/go-kit/kit/endpoint"
	"github.com/micromdm/micromdm/pkg/httputil"
)

func (svc *ConfigService) ApplyVPPToken(ctx context.Context, token_Content []byte) error {

	// Save the data to an SToken
	var sToken = string(token_Content)

	// Create VPPToken with a UDID for ClientContext tracking
	var vppToken VPPToken
	vppToken.UDID = uuid.New().String()
	vppToken.SToken = sToken

	// Convert to JSON
	tokenJSON, err := json.Marshal(vppToken)
	if err != nil {
		return err
	}

	// Save the token
	err = svc.store.AddVPPToken(vppToken.SToken, tokenJSON)
	if err != nil {
		return err
	}
	fmt.Println("stored VPP token", vppToken.SToken)
	return nil
}

type applyVPPTokenRequest struct {
	Content []byte `json:"sToken"`
}

type applyVPPTokenResponse struct {
	Err error `json:"err,omitempty"`
}

func (r applyVPPTokenResponse) Failed() error { return r.Err }

func decodeApplyVPPTokensRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req applyVPPTokenRequest
	err := httputil.DecodeJSONRequest(r, &req)
	return req, err
}

func decodeApplyVPPTokensResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp applyVPPTokenResponse
	err := httputil.DecodeJSONResponse(r, &resp)
	return resp, err
}

func MakeApplyVPPTokensEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(applyVPPTokenRequest)
		err = svc.ApplyVPPToken(ctx, req.Content)
		return applyVPPTokenResponse{
			Err: err,
		}, nil
	}
}

func (e Endpoints) ApplyVPPToken(ctx context.Context, Content []byte) error {
	req := applyVPPTokenRequest{Content: Content}
	resp, err := e.ApplyVPPTokensEndpoint(ctx, req)
	if err != nil {
		return err
	}
	return resp.(applyVPPTokenResponse).Err
}

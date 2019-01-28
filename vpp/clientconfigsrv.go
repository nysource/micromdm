package vpp

import (
	//"encoding/base64"
	"encoding/json"
	"github.com/pkg/errors"
	"strings"
)

// Contains information that associates your particular mdm server to a VPP account token
type ClientContext struct {
	HostName string `json:"hostname"`
	GUID     string `json:"guid"`
}

// Contains location information associated with a VPP account token
type Location struct {
	LocationName string `json:"locationName"`
	LocationID   int    `json:"locationId"`
}

// Contains org information associated with a VPP account token
type ClientConfigSrv struct {
	ClientContext      string   `json:"clientContext"`
	AppleID            string   `json:"appleId,omitempty"`
	OrganizationIDHash string   `json:"organizationIdHash"`
	Status             int      `json:"status"`
	OrganizationID     int      `json:"organizationId"`
	UID                string   `json:"uId"`
	CountryCode        string   `json:"countryCode"`
	Location           Location `json:"location"`
	APNToken           string   `json:"apnToken"`
	Email              string   `json:"email"`
}

// These specify options for the ClientConfigSrv
type GetClientConfigSrvOptions func(*getClientConfigSrvOpts) error

type getClientConfigSrvOpts struct {
	SToken        string `json:"sToken"`
	Verbose       bool   `json:"verbose,omitempty"`
	ClientContext string `json:"clientContext,omitempty"`
}

// Verbose is an optional argument that can be added to GetClientConfigSrv
func VerboseOption(verbose bool) GetClientConfigSrvOptions {
	return func(opts *getClientConfigSrvOpts) error {
		opts.Verbose = verbose
		return nil
	}
}

// ClientContext is an optional argument that can be added to GetClientConfigSrv
func ClientContextOption(context string) GetClientConfigSrvOptions {
	return func(opts *getClientConfigSrvOpts) error {
		opts.ClientContext = context
		return nil
	}
}

// Gets ClientConfigSrv information
func (c *Client) GetClientConfigSrv(token VPPToken, opts ...GetClientConfigSrvOptions) (*ClientConfigSrv, error) {

	// Set required and optional arguments
	request := &getClientConfigSrvOpts{SToken: token.SToken}
	for _, option := range opts {
		if err := option(request); err != nil {
			return nil, err
		}
	}

	// Get the ClientConfigSrvURL
	clientConfigSrvURL := c.VPPServiceConfigSrv.ClientConfigSrvURL

	// Create the ClientConfigSrvURL request
	req, err := c.newRequest("POST", clientConfigSrvURL, request)
	if err != nil {
		return nil, errors.Wrap(err, "create ClientConfigSrv request")
	}

	// Make the request
	var response ClientConfigSrv
	err = c.do(req, &response)

	return &response, errors.Wrap(err, "make ClientConfigSrv request")
}

// Gets the appleID field along with the standard information
func (c *Client) GetClientConfigSrvVerbose(token VPPToken) (*ClientConfigSrv, error) {
	options := VerboseOption(true)
	response, err := c.GetClientConfigSrv(token, options)
	if err != nil {
		return nil, errors.Wrap(err, "using verbose option")
	}

	return response, nil
}

// Configure ClientContext for all VPPTokens
func (c *Client) ConfigureClientContext() (error) {

  // Set Client Context If Needed
	context, err := c.GetClientContext(c.VPPToken)
	if err != nil {
		return errors.Wrap(err, "GetClientContext request")
	}

	if context.HostName != c.ServerPublicURL {
		_, err := c.SetClientContext(c.VPPToken)
		if err != nil {
			return errors.Wrap(err, "SetClientContext request")
		}
	}

	return nil;
}

// Gets the values that determine which mdm server is associated with a VPP account token
func (c *Client) GetClientContext(token VPPToken) (*ClientContext, error) {
	// Get the ClientConfigSrv info
	clientConfigSrv, err := c.GetClientConfigSrv(token)
	if err != nil {
		return nil, errors.Wrap(err, "get ClientContext request")
	}

	// Get the ClientContext string
	var clientContext ClientContext
	var context = clientConfigSrv.ClientContext
	if context != "" {
		// Convert the string to a ClientContext type
		err = json.NewDecoder(strings.NewReader(context)).Decode(&clientContext)
		if err != nil {
			return nil, errors.Wrap(err, "decode ClientContext")
		}
	} else {
		clientContext.HostName = ""
		clientContext.GUID = ""
	}

	return &clientContext, nil
}

// Sets the values that determine which mdm server is associated with a VPP account token
func (c *Client) SetClientContext(token VPPToken) (*ClientContext, error) {

	// Generate a ClientContext string with the new UUID and the current ServerPublicURL
	context := ClientContext{c.ServerPublicURL, token.UDID}
	data, err := json.Marshal(context)
	if err != nil {
		return nil, errors.Wrap(err, "create new ClientContext")
	}
	newContext := string(data)

	// Enter the new ClientContext string into the ClientConfigSrv options
	options := ClientContextOption(newContext)

	// Set the new ClientContext into the VPP account token
	response, err := c.GetClientConfigSrv(token, options)
	if err != nil {
		return nil, errors.Wrap(err, "set ClientContext request")
	}

	// Get the new ClientContext string
	var contextString = response.ClientContext

	// Convert the string to a ClientContext type
	var clientContext ClientContext
	err = json.NewDecoder(strings.NewReader(contextString)).Decode(&clientContext)
	if err != nil {
		return nil, errors.Wrap(err, "decode new ClientContext")
	}

	return &clientContext, nil
}

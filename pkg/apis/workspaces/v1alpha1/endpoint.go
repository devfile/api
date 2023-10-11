//
//
// Copyright Red Hat
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha1

// EndpointProtocol defines the application and transport protocols of the traffic that will go through this endpoint.
// Only one of the following protocols may be specified: http, ws, tcp, udp.
// +kubebuilder:validation:Enum=http;ws;tcp;udp
type EndpointProtocol string

const (
	// Endpoint will have `http` traffic, typically on a TCP connection.
	// It will be automaticaly promoted to `https` when the `secure` field is set to `true`
	HTTPEndpointProtocol EndpointProtocol = "http"
	// Endpoint will have `https` traffic, typically on a TCP connection
	HTTPSEndpointProtocol EndpointProtocol = "https"
	// Endpoint will have `ws` traffic, typically on a TCP connection
	// It will be automaticaly promoted to `wss` when the `secure` field is set to `true`
	WSEndpointProtocol EndpointProtocol = "ws"
	// Endpoint will have `wss` traffic, typically on a TCP connection
	WSSEndpointProtocol EndpointProtocol = "wss"
	// Endpoint will have traffic on a TCP connection,
	// without specifying an application protocol
	TCPEndpointProtocol EndpointProtocol = "tcp"
	// Endpoint will have traffic on an UDP connection,
	// without specifying an application protocol
	UDPEndpointProtocol EndpointProtocol = "udp"
)

// EndpointExposure describes the way an endpoint is exposed on the network.
// Only one of the following exposures may be specified: public, internal, none.
// +kubebuilder:validation:Enum=public;internal;none
type EndpointExposure string

const (
	// Endpoint will be exposed on the public network, typically through
	// a K8S ingress or an OpenShift route
	PublicEndpointExposure EndpointExposure = "public"
	// Endpoint will be exposed internally outside of the main workspace POD,
	// typically by K8S services, to be consumed by other elements running
	// on the same cloud internal network.
	InternalEndpointExposure EndpointExposure = "internal"
	// Endpoint will not be exposed and will only be accessible
	// inside the main workspace POD, on a local address.
	NoneEndpointExposure EndpointExposure = "none"
)

type Endpoint struct {
	Name string `json:"name"`

	// +optional
	TargetPort int `json:"targetPort,omitempty"`

	// Describes how the endpoint should be exposed on the network.
	//
	// - `public` means that the endpoint will be exposed on the public network, typically through
	// a K8S ingress or an OpenShift route.
	//
	// - `internal` means that the endpoint will be exposed internally outside of the main workspace POD,
	// typically by K8S services, to be consumed by other elements running
	// on the same cloud internal network.
	//
	// - `none` means that the endpoint will not be exposed and will only be accessible
	// inside the main workspace POD, on a local address.
	//
	// Default value is `public`
	// +optional
	Exposure EndpointExposure `json:"exposure,omitempty"`

	// Describes the application and transport protocols of the traffic that will go through this endpoint.
	//
	// - `http`: Endpoint will have `http` traffic, typically on a TCP connection.
	// It will be automaticaly promoted to `https` when the `secure` field is set to `true`.
	//
	// - `https`: Endpoint will have `https` traffic, typically on a TCP connection.
	//
	// - `ws`: Endpoint will have `ws` traffic, typically on a TCP connection.
	// It will be automaticaly promoted to `wss` when the `secure` field is set to `true`.
	//
	// - `wss`: Endpoint will have `wss` traffic, typically on a TCP connection.
	//
	// - `tcp`: Endpoint will have traffic on a TCP connection, without specifying an application protocol.
	//
	// - `udp`: Endpoint will have traffic on an UDP connection, without specifying an application protocol.
	//
	// Default value is `http`
	// +optional
	Protocol string `json:"protocol,omitempty"`

	// Describes whether the endpoint should be secured and protected by some
	// authentication process
	// +optional
	Secure bool `json:"secure,omitempty"`

	// Path of the endpoint URL
	// +optional
	Path string `json:"path,omitempty"`

	// Map of implementation-dependant string-based free-form attributes.
	//
	// Examples of Che-specific attributes:
	//
	// - cookiesAuthEnabled: "true" / "false",
	//
	// - type: "terminal" / "ide" / "ide-dev",
	// +optional
	Attributes map[string]string `json:"attributes,omitempty"`
}

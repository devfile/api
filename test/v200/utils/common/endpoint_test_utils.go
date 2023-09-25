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

package common

import (
	"fmt"

	schema "github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
)

var Exposures = [...]schema.EndpointExposure{schema.PublicEndpointExposure, schema.InternalEndpointExposure, schema.NoneEndpointExposure}

// getRandomExposure returns a random exposure value
func getRandomExposure() schema.EndpointExposure {
	return Exposures[GetRandomNumber(1, len(Exposures))-1]
}

var Protocols = [...]schema.EndpointProtocol{schema.HTTPEndpointProtocol, schema.HTTPSEndpointProtocol, schema.WSEndpointProtocol, schema.WSSEndpointProtocol, schema.TCPEndpointProtocol, schema.UDPEndpointProtocol}

// getRandomProtocol returns a random protocol value
func getRandomProtocol() schema.EndpointProtocol {
	return Protocols[GetRandomNumber(1, len(Protocols))-1]
}

// getUniquePort return a port value not previously used in that same devfile
func (devfile *TestDevfile) getUniquePort() int {

	// max sure a lot of unique ports exist
	maxPorts := len(devfile.UsedPorts) + 5000

	var port int
	used := true
	for used {
		port = GetRandomNumber(1, maxPorts)
		_, used = devfile.UsedPorts[port]
	}
	devfile.UsedPorts[port] = true
	return port
}

// CreateEndpoints creates and returns a randon number of endpoints in a schema structure
func (devfile *TestDevfile) CreateEndpoints() []schema.Endpoint {

	numEndpoints := GetRandomNumber(1, 5)
	endpoints := make([]schema.Endpoint, numEndpoints)

	commonPort := devfile.getUniquePort()

	for i := 0; i < numEndpoints; i++ {

		endpoint := schema.Endpoint{}

		endpoint.Name = GetRandomUniqueString(GetRandomNumber(5, 10), true)
		LogInfoMessage(fmt.Sprintf("   ....... add endpoint %d name  : %s", i, endpoint.Name))

		if GetBinaryDecision() {
			endpoint.TargetPort = devfile.getUniquePort()
		} else {
			endpoint.TargetPort = commonPort
		}
		LogInfoMessage(fmt.Sprintf("   ....... add endpoint %d targetPort: %d", i, endpoint.TargetPort))

		if GetBinaryDecision() {
			endpoint.Exposure = getRandomExposure()
			LogInfoMessage(fmt.Sprintf("   ....... add endpoint %d exposure: %s", i, endpoint.Exposure))
		}

		if GetBinaryDecision() {
			endpoint.Protocol = getRandomProtocol()
			LogInfoMessage(fmt.Sprintf("   ....... add endpoint %d protocol: %s", i, endpoint.Protocol))
		}

		value := GetBinaryDecision()
		endpoint.Secure = &value
		LogInfoMessage(fmt.Sprintf("   ....... add endpoint %d secure: %t", i, *endpoint.Secure))

		if GetBinaryDecision() {
			endpoint.Path = "/Path_" + GetRandomString(GetRandomNumber(3, 15), false)
			LogInfoMessage(fmt.Sprintf("   ....... add endpoint %d path: %s", i, endpoint.Path))
		}

		endpoints[i] = endpoint

	}

	return endpoints
}

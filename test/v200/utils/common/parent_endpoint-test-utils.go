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

// CreateParentEndpoints creates and returns a random number of endpoints in a schema structure
func (devfile *TestDevfile) CreateParentEndpoints() []schema.EndpointParentOverride {

	numEndpoints := GetRandomNumber(1, 5)
	endpoints := make([]schema.EndpointParentOverride, numEndpoints)

	commonPort := devfile.getUniquePort()

	for i := 0; i < numEndpoints; i++ {

		endpoint := schema.EndpointParentOverride{}

		endpoint.Name = GetRandomUniqueString(GetRandomNumber(5, 10), true)
		LogInfoMessage(fmt.Sprintf("   ....... add endpoint %d name  : %s", i, endpoint.Name))

		if GetBinaryDecision() {
			endpoint.TargetPort = devfile.getUniquePort()
		} else {
			endpoint.TargetPort = commonPort
		}
		LogInfoMessage(fmt.Sprintf("   ....... add endpoint %d targetPort: %d", i, endpoint.TargetPort))

		if GetBinaryDecision() {
			endpoint.Exposure = schema.EndpointExposureParentOverride(getRandomExposure())
			LogInfoMessage(fmt.Sprintf("   ....... add endpoint %d exposure: %s", i, endpoint.Exposure))
		}

		if GetBinaryDecision() {
			endpoint.Protocol = schema.EndpointProtocolParentOverride(getRandomProtocol())
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

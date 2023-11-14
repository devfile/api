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

	header "github.com/devfile/api/v2/pkg/devfile"
	"github.com/lucasjones/reggen"
)

// AddMetadata populates the metadata object with random attributes
func (devfile *TestDevfile) AddMetaData() header.DevfileMetadata {
	LogInfoMessage("metadata added")
	//point to the intialized struct
	metadata := &devfile.SchemaDevFile.DevfileHeader.Metadata
	devfile.MetaDataAdded(metadata)
	devfile.setMetaDataValues(metadata)
	return *metadata
}

// MetadataAdded notifies a registered follower
func (devfile *TestDevfile) MetaDataAdded(metadata *header.DevfileMetadata) {
	if devfile.Follower != nil {
		devfile.Follower.SetMetaData(*metadata)
	}

}

// MetaDataUpdated notifies a registered follower that the metadata has been updated
func (devfile *TestDevfile) MetaDataUpdated(metadata *header.DevfileMetadata) {
	LogInfoMessage("metadata updated")
	if devfile.Follower != nil {
		devfile.Follower.UpdateMetaData(*metadata)
	}
}

// setProperty randomly sets the string properties of the metadata object.
func setProperty(propertyName string, property *string) {
	if GetRandomDecision(2, 1) {
		*property = GetRandomString(8, false)
		LogInfoMessage(fmt.Sprintf("   ....... %s %s", propertyName, *property))
	}
}

// setMetadataValues randomly adds/modifies metadata object properties.  Since these are optional properties, the test is
// set up so they are twice as likely to appear in the generated files to ensure sufficient coverage.
func (devfile *TestDevfile) setMetaDataValues(metadata *header.DevfileMetadata) {
	setProperty("Description", &metadata.Description)
	setProperty("DisplayName", &metadata.DisplayName)
	setProperty("GlobalMemoryLimit", &metadata.GlobalMemoryLimit)
	setProperty("Icon", &metadata.Icon)
	setProperty("Name", &metadata.Name)

	if GetRandomDecision(2, 1) {

		numTags := GetRandomNumber(1, 5)
		LogInfoMessage(fmt.Sprintf("   ....... add %d tag(s) to tags ", numTags))
		for i := 0; i < numTags; i++ {
			metadata.Tags = append(metadata.Tags, GetRandomString(8, false))
		}

	}

	if GetRandomDecision(2, 1) {
		//generate a valid version string based on the regex that's in the spec.  Limit each segment to max 3 characters. e.g. 339.957.11-t.9+-.nkJ will be generated
		version, err := reggen.Generate("^([0-9]+)\\.([0-9]+)\\.([0-9]+)(\\-[0-9a-z-]+(\\.[0-9a-z-]+)*)?(\\+[0-9A-Za-z-]+(\\.[0-9A-Za-z-]+)*)?$", 3)
		if err != nil {
			LogErrorMessage("Failed to generate random version, schema version will be used")
			version = devfile.SchemaDevFile.SchemaVersion
		}

		metadata.Version = version
		LogInfoMessage(fmt.Sprintf("   ....... Version %s", metadata.Version))
	}

	devfile.MetaDataUpdated(metadata)
}

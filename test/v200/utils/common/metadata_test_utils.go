package common

import (
	"fmt"

	header "github.com/devfile/api/v2/pkg/devfile"
	"github.com/lucasjones/reggen"
)

func (devfile *TestDevfile) AddMetaData() header.DevfileMetadata {
	LogInfoMessage("metadata added")
	//point to the intialized struct
	metadata := &devfile.SchemaDevFile.DevfileHeader.Metadata
	devfile.MetaDataAdded(metadata)
	devfile.setMetaDataValues(metadata)
	return *metadata
}

// MetadataAdded adds metadata to the test schema and notifies a registered follower
func (devfile *TestDevfile) MetaDataAdded(metadata *header.DevfileMetadata) {
	if devfile.Follower != nil {
		devfile.Follower.SetMetaData(*metadata)
	}

}

func (devfile *TestDevfile) MetaDataUpdated(metadata *header.DevfileMetadata) {
	LogInfoMessage("metadata updated")
	if devfile.Follower != nil {
		devfile.Follower.UpdateMetaData(*metadata)
	}
}

func setProperty(propertyName string, property *string) {
	if GetRandomDecision(2, 1) {
		*property = GetRandomString(8, false)
		LogInfoMessage(fmt.Sprintf("   ....... %s %s", propertyName, *property))
	}
}

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

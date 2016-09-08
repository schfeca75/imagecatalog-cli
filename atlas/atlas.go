package atlas

import (
	"fmt"
	"github.com/hashicorp/atlas-go/v1"
	s "strings"
)

type CloudProvider interface {
	GetIdentifier() string
	CreateImageMap(artifactVersion *atlas.ArtifactVersion) map[string]string
}

type CloudProviderBase struct {
	Identifier string
}

func (cp CloudProviderBase) GetIdentifier() string {
	return cp.Identifier
}

func (cp CloudProviderBase) CreateImageMap(artifactVersion *atlas.ArtifactVersion) map[string]string {
	return nil
}

type Amazon struct {
	CloudProvider
}

func (Amazon) CreateImageMap(artifactVersion *atlas.ArtifactVersion) map[string]string {
	metaData := artifactVersion.Metadata
	imageMap := make(map[string]string)
	for k, v := range metaData {
		if s.HasPrefix(k, "region.") {
			imageMap[s.TrimPrefix(k, "region.")] = v
		}
	}
	return imageMap
}

type Azure struct {
	CloudProvider
}

func (Azure) CreateImageMap(artifactVersion *atlas.ArtifactVersion) map[string]string {
	metaData := artifactVersion.Metadata
	imageName := metaData["image_name"]
	imageMap := map[string]string{
		"East Asia":        "https://sequenceiqeastasia2.blob.core.windows.net/images/" + imageName + ".vhd",
		"East US":          "https://sequenceiqeastus12.blob.core.windows.net/images/" + imageName + ".vhd",
		"Central US":       "https://sequenceiqcentralus2.blob.core.windows.net/images/" + imageName + ".vhd",
		"North Europe":     "https://sequenceiqnortheurope2.blob.core.windows.net/images/" + imageName + ".vhd",
		"South Central US": "https://sequenceiqouthcentralus2.blob.core.windows.net/images/" + imageName + ".vhd",
		"North Central US": "https://sequenceiqorthcentralus2.blob.core.windows.net/images/" + imageName + ".vhd",
		"East US 2":        "https://sequenceiqeastus22.blob.core.windows.net/images/" + imageName + ".vhd",
		"Japan East":       "https://sequenceiqjapaneast2.blob.core.windows.net/images/" + imageName + ".vhd",
		"Japan West":       "https://sequenceiqjapanwest2.blob.core.windows.net/images/" + imageName + ".vhd",
		"Southeast Asia":   "https://sequenceiqsoutheastasia2.blob.core.windows.net/images/" + imageName + ".vhd",
		"West US":          "https://sequenceiqwestus2.blob.core.windows.net/images/" + imageName + ".vhd",
		"West Europe":      "https://sequenceiqwesteurope2.blob.core.windows.net/images/" + imageName + ".vhd",
		"Brazil South":     "https://sequenceiqbrazilsouth2.blob.core.windows.net/images/" + imageName + ".vhd",
		"Canada East":      "https://sequenceiqcanadaeast.blob.core.windows.net/images/" + imageName + ".vhd",
		"Canada Central":   "https://sequenceiqcanadacentral.blob.core.windows.net/images/" + imageName + ".vhd",
	}
	return imageMap
}

type Gcp struct {
	CloudProvider
}

func (Gcp) CreateImageMap(artifactVersion *atlas.ArtifactVersion) map[string]string {
	imageMap := map[string]string{
		"default": "sequenceiqimage/" + artifactVersion.ID + ".tar.gz",
	}
	return imageMap
}

type Openstack struct {
	CloudProvider
}

func (Openstack) CreateImageMap(artifactVersion *atlas.ArtifactVersion) map[string]string {
	imageMap := map[string]string{
		"default": s.TrimSuffix(artifactVersion.ID, ".img"),
	}
	return imageMap
}

var cloudProviders = map[string]CloudProvider{
	"amazon":        Amazon{CloudProviderBase{"aws"}},
	"azure":         Azure{CloudProviderBase{"azure_rm"}},
	"googlecompute": Gcp{CloudProviderBase{"gcp"}},
	"openstack":     Openstack{CloudProviderBase{"openstack"}},
}

func GetCloudImages(ambariVersion string, hdpVersion string) (map[string]map[string]string, error) {
	client := atlas.DefaultClient()
	imageMap := make(map[string]map[string]string)
	for k, v := range cloudProviders {
		searchOpts := atlas.ArtifactSearchOpts{
			User: "sequenceiq", Name: "cloudbreak", Type: k + ".image",
			Metadata: map[string]string{"ambari_version": ambariVersion, "hdp_version": hdpVersion}}
		artifactVersions, err := client.ArtifactSearch(&searchOpts)
		if err != nil {
			return nil, err // TODO provider information on error
		}
		fmt.Printf("%#v\n", artifactVersions)
		if len(artifactVersions) > 0 {
			imageMap[v.GetIdentifier()] = v.CreateImageMap(artifactVersions[0]) // Select the latest version
		}
	}
	return imageMap, nil
}

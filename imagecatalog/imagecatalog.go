package imagecatalog

import (
	"encoding/json"
	"io/ioutil"
)

type HdpRepository struct {
	StackRepository map[string]string `json:"stack"`
	UtilRepository  map[string]string `json:"util"`
}

type HdpCatalog struct {
	Version    string                       `json:"version"`
	Repository HdpRepository                `json:"repo"`
	Images     map[string]map[string]string `json:"images"`
}

type AmbariCatalog struct {
	CbVersions []string          `json:"cb_versions"`
	Version    string            `json:"version"`
	Repository map[string]string `json:"repo"`
	HdpCatalog []HdpCatalog      `json:"hdp"`
}

type CloudbreakCatalog struct {
	AmbariCatalog AmbariCatalog `json:"ambari"`
}

type ImageCatalog struct {
	CloudbreakCatalog []CloudbreakCatalog `json:"cloudbreak"`
}

type CbImageInfo struct {
	CloudbreakVersions map[string]bool
	AmbariRepository   map[string]string
	ImageInfoMap       map[string]ImageInfo
}

type ImageInfo struct {
	HdpRepository       HdpRepository
	CloudProviderImages map[string]map[string]string
}

func ParseImageCatalog(fileName string) (*ImageCatalog, error) {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var imageCatalog ImageCatalog
	err = json.Unmarshal(file, &imageCatalog)
	if err != nil {
		return nil, err
	}
	return &imageCatalog, nil
}

func CreateVersionInfoMap(imageCatalog *ImageCatalog) map[string]CbImageInfo {
	versionInfoMap := make(map[string]CbImageInfo)
	for _, ambari := range imageCatalog.CloudbreakCatalog {
		hdpVersionMap := make(map[string]ImageInfo)
		cbVersions := make(map[string]bool)
		for _, cbVersion := range ambari.AmbariCatalog.CbVersions {
			cbVersions[cbVersion] = true
		}
		for _, hdp := range ambari.AmbariCatalog.HdpCatalog {
			imageInfo := ImageInfo{
				CloudProviderImages: hdp.Images,
				HdpRepository:       hdp.Repository}
			hdpVersionMap[hdp.Version] = imageInfo
		}
		cbImageInfo := CbImageInfo{
			CloudbreakVersions: cbVersions,
			AmbariRepository:   ambari.AmbariCatalog.Repository,
			ImageInfoMap:       hdpVersionMap}
		versionInfoMap[ambari.AmbariCatalog.Version] = cbImageInfo
	}
	return versionInfoMap
}

func CreateImageCatalog(cbImageInfoMap map[string]CbImageInfo) (*ImageCatalog, error) {
	imageCatalog := ImageCatalog{CloudbreakCatalog: []CloudbreakCatalog{}}
	sortedAmbariVersions := SortCbVersionKeys(cbImageInfoMap)
	for _, ambariVersion := range sortedAmbariVersions {
		cbImageInfo := cbImageInfoMap[ambariVersion]
		ambari := AmbariCatalog{CbVersions: []string{}, Version: ambariVersion, HdpCatalog: []HdpCatalog{}} // Repo information is missing currently
		for cbVersion := range cbImageInfo.CloudbreakVersions {
			ambari.CbVersions = append(ambari.CbVersions, cbVersion)
		}
		SortVersions(ambari.CbVersions)
		sortedHdpVersions := SortImVersionKeys(cbImageInfo.ImageInfoMap)
		for _, hdpVersion := range sortedHdpVersions {
			imageInfo := cbImageInfo.ImageInfoMap[hdpVersion]
			hdp := HdpCatalog{Version: hdpVersion, Images: imageInfo.CloudProviderImages} // missing hdprepo, need to be stored in ImageInfo
			ambari.HdpCatalog = append(ambari.HdpCatalog, hdp)
		}
		imageCatalog.CloudbreakCatalog = append(imageCatalog.CloudbreakCatalog, CloudbreakCatalog{AmbariCatalog: ambari})
	}
	return &imageCatalog, nil
}

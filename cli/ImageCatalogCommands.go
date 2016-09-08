package cli

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/hortonworks/imagecatalog-cli/atlas"
	ic "github.com/hortonworks/imagecatalog-cli/imagecatalog"
	"github.com/urfave/cli"
	"io/ioutil"
)

var (
	FlDebug = cli.BoolFlag{
		Name:   "debug",
		Usage:  "debug mode",
		EnvVar: "DEBUG",
	}
	FlImageCatalog = cli.StringFlag{
		Name:   "imageCatalog",
		Usage:  "the image catalog file",
		EnvVar: "IMAGE_CATALOG",
	}
	FlOutputImageCatalog = cli.StringFlag{
		Name:   "outputImageCatalog",
		Usage:  "the output image catalog file",
		EnvVar: "OUTPUT_IMAGE_CATALOG",
	}
	FlCloudbreakVersion = cli.StringFlag{
		Name:   "cloudbreakVersion",
		Usage:  "cloudbreak version",
		EnvVar: "CLOUDBREAK_VERSION",
	}
	FlAmbariVersion = cli.StringFlag{
		Name:   "ambariVersion",
		Usage:  "ambariVersion",
		EnvVar: "AMBARI_VERSION",
	}
	FlHdpVersion = cli.StringFlag{
		Name:   "hdpVersion",
		Usage:  "hdpVersion",
		EnvVar: "HDP_VERSION",
	}
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func validateFields(c *cli.Context) {
	if len(c.String(FlImageCatalog.Name)) == 0 {
		ExitWithErrorMsg(FlImageCatalog.Name + " is mandatory")
	}
	if len(c.String(FlCloudbreakVersion.Name)) == 0 {
		ExitWithErrorMsg(FlCloudbreakVersion.Name + " is mandatory")
	}
	if len(c.String(FlAmbariVersion.Name)) == 0 {
		ExitWithErrorMsg(FlAmbariVersion.Name + " is mandatory")
	}
	if len(c.String(FlHdpVersion.Name)) == 0 {
		ExitWithErrorMsg(FlHdpVersion.Name + " is mandatory")
	}
}

func AddVersion(c *cli.Context) {
	validateFields(c)

	imageCatalogFile := c.String(FlImageCatalog.Name)
	outputImageCatalogFile := c.String(FlOutputImageCatalog.Name)
	if len(outputImageCatalogFile) == 0 {
		outputImageCatalogFile = imageCatalogFile
	}
	cbVersion := c.String(FlCloudbreakVersion.Name)
	ambariVersion := c.String(FlAmbariVersion.Name)
	hdpVersion := c.String(FlHdpVersion.Name)

	log.Infof("[AddVersion] Adding new version [cb: %s / ambar: %s / hdp: %s], input: %s, output: %s", cbVersion, ambariVersion, hdpVersion, imageCatalogFile,
		outputImageCatalogFile)

	imageCatalog, err := ic.ParseImageCatalog(imageCatalogFile)
	ExitOnError(err, "Parsing imagecatalog was unsuccessful")
	log.Debugf("ImageCatalog: %+v\n", *imageCatalog)

	versionInfoMap := ic.CreateVersionInfoMap(imageCatalog)
	log.Debugf("VersionInfoMap created from ImageCatalog: %+v\n", versionInfoMap)

	cbImageInfo, ambariVersionFound := versionInfoMap[ambariVersion]
	if ambariVersionFound {
		if !cbImageInfo.CloudbreakVersions[cbVersion] {
			cbImageInfo.CloudbreakVersions[cbVersion] = true
		}
	} else {
		cbImageInfo = ic.CbImageInfo{
			CloudbreakVersions: map[string]bool{cbVersion: true},
			ImageInfoMap:       map[string]ic.ImageInfo{}}
		versionInfoMap[ambariVersion] = cbImageInfo
	}

	cloudImages, err := atlas.GetCloudImages(ambariVersion, hdpVersion)
	ExitOnError(err, "Cannot get artifact information from atlas for cloudprovider")
	cbImageInfo.ImageInfoMap[hdpVersion] = ic.ImageInfo{CloudProviderImages: cloudImages}

	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Printf("%+v\n", versionInfoMap[ambariVersion].ImageInfoMap[hdpVersion])
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Printf("%+v\n", versionInfoMap)
	for _, v := range versionInfoMap {
		fmt.Println()
		fmt.Println()
		fmt.Println()
		fmt.Printf("%+v\n", v)
	}
	imageCatalog2, _ := ic.CreateImageCatalog(versionInfoMap)
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Printf("%+v\n", imageCatalog2)
	json, _ := json.MarshalIndent(imageCatalog2, "", "  ")
	ioutil.WriteFile(outputImageCatalogFile, json, 0644)
}

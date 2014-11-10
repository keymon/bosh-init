package config

import (
	"encoding/json"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	boshsys "github.com/cloudfoundry/bosh-agent/system"
)

type fileSystemDeploymentConfigService struct {
	configPath string
	fs         boshsys.FileSystem
	logger     boshlog.Logger
	logTag     string
}

func NewFileSystemDeploymentConfigService(configPath string, fs boshsys.FileSystem, logger boshlog.Logger) DeploymentConfigService {
	return fileSystemDeploymentConfigService{
		configPath: configPath,
		fs:         fs,
		logger:     logger,
		logTag:     "config",
	}
}

type StemcellRecord struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	SHA1    string `json:"sha1"`
	CID     string `json:"cid"`
}

type DeploymentFile struct {
	UUID      string           `json:"uuid"`
	Stemcells []StemcellRecord `json:"stemcells"`
	VMCID     string           `json:"vm_cid"`
	DiskCID   string           `json:"disk_cid"`
}

func (s fileSystemDeploymentConfigService) Load() (DeploymentConfig, error) {
	config := DeploymentConfig{}

	if !s.fs.FileExists(s.configPath) {
		return config, nil
	}

	deploymentFileContents, err := s.fs.ReadFile(s.configPath)
	if err != nil {
		return config, bosherr.WrapError(err, "Reading deployment config file `%s'", s.configPath)
	}
	s.logger.Debug(s.logTag, "Deployment File Contents %#s", deploymentFileContents)

	deploymentFile := DeploymentFile{}

	err = json.Unmarshal(deploymentFileContents, &deploymentFile)
	if err != nil {
		return config, bosherr.WrapError(err, "Unmarshalling deployment config file `%s'", s.configPath)
	}

	config.DeploymentUUID = deploymentFile.UUID
	config.Stemcells = deploymentFile.Stemcells
	config.VMCID = deploymentFile.VMCID
	config.DiskCID = deploymentFile.DiskCID

	return config, nil
}

func (s fileSystemDeploymentConfigService) Save(config DeploymentConfig) error {
	deploymentFile := DeploymentFile{
		UUID:      config.DeploymentUUID,
		Stemcells: config.Stemcells,
		VMCID:     config.VMCID,
		DiskCID:   config.DiskCID,
	}
	jsonContent, err := json.MarshalIndent(deploymentFile, "", "    ")
	if err != nil {
		return bosherr.WrapError(err, "Marshalling deployment config into JSON")
	}

	err = s.fs.WriteFile(s.configPath, jsonContent)
	if err != nil {
		return bosherr.WrapError(err, "Writing deployment config file `%s'", s.configPath)
	}

	return nil
}

package agentclient

import (
	bmas "github.com/cloudfoundry/bosh-micro-cli/deployment/applyspec"
)

type AgentClient interface {
	Ping() (string, error)
	Stop() error
	Apply(bmas.ApplySpec) error
	Start() error
	GetState() (AgentState, error)
	MountDisk(string) error
	UnmountDisk(string) error
	ListDisk() ([]string, error)
	MigrateDisk() error
	CompilePackage(packageSource BlobRef, compiledPackageDependencies []BlobRef) (compiledPackageRef BlobRef, err error)
}

type AgentState struct {
	JobState string
}

type BlobRef struct {
	Name        string
	Version     string
	BlobstoreID string
	SHA1        string
}
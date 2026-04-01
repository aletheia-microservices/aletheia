package components

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type ServiceInfo struct {
	Name            string
	Filepath        string
	Package         string
	PackagePath     string
	ConstructorName string
	ServiceArgs     []string
	Methods         []string
	Edges           []string
}

func (info *ServiceInfo) String() string {
	return fmt.Sprintf("Name: %s, Filepath: %s, ConstructorName: %s, ServiceArgs: %v, Methods: %v", info.Name, info.Filepath, info.ConstructorName, info.ServiceArgs, info.Methods)
}

func (info *ServiceInfo) GetArgumentAt(idx int) string {
	if idx < len(info.ServiceArgs) {
		return info.ServiceArgs[idx]
	}
	logrus.Fatalf("[SERVICE INFO] [%s] index %d out of bounds for service info with %d arguments: %v", info.Name, idx, len(info.ServiceArgs), info.ServiceArgs)
	return ""
}

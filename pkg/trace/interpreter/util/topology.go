package util

import (
	"fmt"
)

// ServiceTypeName returns the default service type
const ServiceTypeName = "service"

// ServiceInstanceTypeName returns the default service instance type
const ServiceInstanceTypeName = "service-instance"

// CreateServiceURN creates the urn identifier for all service components
func CreateServiceURN(serviceName string) string {
	return fmt.Sprintf("urn:%s:/%s", ServiceTypeName, serviceName)
}

// CreateServiceInstanceURN creates the urn identifier for all service instance components
func CreateServiceInstanceURN(serviceName string, hostname string, pid int, createTime int64) string {
	return fmt.Sprintf("urn:%s:/%s:/%s:%d:%d", ServiceInstanceTypeName, serviceName, hostname, pid, createTime)
}

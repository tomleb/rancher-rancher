/*
Copyright 2024 Rancher Labs, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by main. DO NOT EDIT.

package v3

import (
	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"github.com/rancher/wrangler/v3/pkg/generic"
)

// ProjectNetworkPolicyController interface for managing ProjectNetworkPolicy resources.
type ProjectNetworkPolicyController interface {
	generic.ControllerInterface[*v3.ProjectNetworkPolicy, *v3.ProjectNetworkPolicyList]
}

type ProjectNetworkPolicyControllerContext interface {
	generic.ControllerInterfaceContext[*v3.ProjectNetworkPolicy, *v3.ProjectNetworkPolicyList]
}

// ProjectNetworkPolicyClient interface for managing ProjectNetworkPolicy resources in Kubernetes.
type ProjectNetworkPolicyClient interface {
	generic.ClientInterface[*v3.ProjectNetworkPolicy, *v3.ProjectNetworkPolicyList]
}

// ProjectNetworkPolicyCache interface for retrieving ProjectNetworkPolicy resources in memory.
type ProjectNetworkPolicyCache interface {
	generic.CacheInterface[*v3.ProjectNetworkPolicy]
}

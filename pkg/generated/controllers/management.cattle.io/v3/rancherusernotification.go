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

// RancherUserNotificationController interface for managing RancherUserNotification resources.
type RancherUserNotificationController interface {
	generic.NonNamespacedControllerInterface[*v3.RancherUserNotification, *v3.RancherUserNotificationList]
}

type RancherUserNotificationControllerContext interface {
	generic.NonNamespacedControllerInterfaceContext[*v3.RancherUserNotification, *v3.RancherUserNotificationList]
}

// RancherUserNotificationClient interface for managing RancherUserNotification resources in Kubernetes.
type RancherUserNotificationClient interface {
	generic.NonNamespacedClientInterface[*v3.RancherUserNotification, *v3.RancherUserNotificationList]
}

// RancherUserNotificationCache interface for retrieving RancherUserNotification resources in memory.
type RancherUserNotificationCache interface {
	generic.NonNamespacedCacheInterface[*v3.RancherUserNotification]
}

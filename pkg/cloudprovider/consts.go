// Copyright 2019 Yunion
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cloudprovider

import (
	"yunion.io/x/pkg/errors"
)

const (
	CloudVMStatusRunning      = "running"
	CloudVMStatusSuspend      = "suspend"
	CloudVMStatusStopped      = "stopped"
	CloudVMStatusChangeFlavor = "change_flavor"
	CloudVMStatusDeploying    = "deploying"
	CloudVMStatusOther        = "other"

	ErrNotFound            = errors.Error("id not found")
	ErrDuplicateId         = errors.Error("duplicate id")
	ErrInvalidStatus       = errors.Error("invalid status")
	ErrTimeout             = errors.Error("timeout")
	ErrNotImplemented      = errors.Error("Not implemented")
	ErrNotSupported        = errors.Error("Not supported")
	ErrInvalidProvider     = errors.Error("Invalid provider")
	ErrNoBalancePermission = errors.Error("No balance permission")
)

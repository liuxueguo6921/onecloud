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

package qcloud

import (
	"fmt"
	"time"

	api "yunion.io/x/onecloud/pkg/apis/compute"
	"yunion.io/x/onecloud/pkg/multicloud"

	"github.com/pkg/errors"
	"yunion.io/x/onecloud/pkg/cloudprovider"
)

type SPrivateIpAddress struct {
	nic              *SNetworkInterface
	Description      string
	Primary          bool
	PrivateIpAddress string
	PublicIpAddress  string
	IsWanIpBlocked   bool
	State            string
}

func (ip *SPrivateIpAddress) GetGlobalId() string {
	return ip.PrivateIpAddress
}

func (ip *SPrivateIpAddress) GetIP() string {
	return ip.PrivateIpAddress
}

func (ip *SPrivateIpAddress) GetINetworkId() string {
	return ip.nic.SubnetId
}

func (ip *SPrivateIpAddress) IsPrimary() bool {
	return ip.Primary
}

type SNetworkInterface struct {
	multicloud.SNetworkInterfaceBase
	region                      *SRegion
	VpcId                       string
	SubnetId                    string
	NetworkInterfaceId          string
	NetworkInterfaceName        string
	NetworkInterfaceDescription string
	GroupSet                    []string
	Primary                     bool
	MacAddress                  string
	State                       string
	CreatedTime                 time.Time
	Attachment                  string
	Zone                        string
	PrivateIpAddressSet         []SPrivateIpAddress
}

func (nic *SNetworkInterface) GetName() string {
	return nic.NetworkInterfaceName
}

func (nic *SNetworkInterface) GetId() string {
	return nic.NetworkInterfaceId
}

func (nic *SNetworkInterface) GetGlobalId() string {
	return nic.NetworkInterfaceId
}

func (nic *SNetworkInterface) GetMacAddress() string {
	return nic.MacAddress
}

func (nic *SNetworkInterface) GetAssociateType() string {
	return api.NETWORK_INTERFACE_ASSOCIATE_TYPE_SERVER
}

func (nic *SNetworkInterface) GetAssociateId() string {
	return nic.Attachment
}

func (nic *SNetworkInterface) GetStatus() string {
	switch nic.State {
	case "PENDING":
		return api.NETWORK_INTERFACE_STATUS_CREATING
	case "AVAILABLE":
		return api.NETWORK_INTERFACE_STATUS_AVAILABLE
	case "ATTACHING":
		return api.NETWORK_INTERFACE_STATUS_ATTACHING
	case "DETACHING":
		return api.NETWORK_INTERFACE_STATUS_DETACHING
	case "DELETING":
		return api.NETWORK_INTERFACE_STATUS_DELETING
	}
	return nic.State
}

func (nic *SNetworkInterface) GetICloudInterfaceAddresses() ([]cloudprovider.ICloudInterfaceAddress, error) {
	address := []cloudprovider.ICloudInterfaceAddress{}
	for i := 0; i < len(nic.PrivateIpAddressSet); i++ {
		nic.PrivateIpAddressSet[i].nic = nic
		address = append(address, &nic.PrivateIpAddressSet[i])
	}
	return address, nil
}

func (region *SRegion) GetINetworkInterfaces() ([]cloudprovider.ICloudNetworkInterface, error) {
	interfaces := []SNetworkInterface{}
	for {
		parts, total, err := region.GetNetworkInterfaces([]string{}, "", len(interfaces), 50)
		if err != nil {
			return nil, err
		}
		interfaces = append(interfaces, parts...)
		if len(interfaces) >= total {
			break
		}
	}
	ret := []cloudprovider.ICloudNetworkInterface{}
	for i := 0; i < len(interfaces); i++ {
		interfaces[i].region = region
		ret = append(ret, &interfaces[i])
	}
	return ret, nil
}

func (region *SRegion) GetNetworkInterfaces(interfaceIds []string, subnetId string, offset int, limit int) ([]SNetworkInterface, int, error) {
	if limit > 50 || limit <= 0 {
		limit = 50
	}
	params := map[string]string{}
	params["Limit"] = fmt.Sprintf("%d", limit)
	params["Offset"] = fmt.Sprintf("%d", offset)

	for idx, interfaceId := range interfaceIds {
		params[fmt.Sprintf("NetworkInterfaceIds.%d", idx)] = interfaceId
	}

	if len(subnetId) > 0 {
		params["Filters.0.Name"] = "subnet-id"
		params["Filters.0.Values.0"] = subnetId
	}
	body, err := region.vpcRequest("DescribeNetworkInterfaces", params)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "DescribeNetworkInterfaces")
	}

	interfaces := []SNetworkInterface{}
	err = body.Unmarshal(&interfaces, "NetworkInterfaceSet")
	if err != nil {
		return nil, 0, errors.Wrapf(err, "Unmarshal.NetworkInterfaceSet")
	}
	total, _ := body.Float("TotalCount")
	return interfaces, int(total), nil
}

func (region *SRegion) DeleteNetworkInterface(interfaceId string) error {
	params := map[string]string{}
	params["Region"] = region.Region
	params["NetworkInterfaceId"] = interfaceId

	_, err := region.vpcRequest("DeleteNetworkInterface", params)
	if err != nil {
		return errors.Wrapf(err, "vpcRequest.DeleteNetworkInterface")
	}
	return nil
}

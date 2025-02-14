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

package huawei

import (
	api "yunion.io/x/onecloud/pkg/apis/compute"
	"yunion.io/x/onecloud/pkg/cloudprovider"
	"yunion.io/x/onecloud/pkg/multicloud"
	"yunion.io/x/pkg/utils"
)

type SFixedIP struct {
	IpAddress string
	SubnetID  string
	NetworkId string
}

func (fixip *SFixedIP) GetGlobalId() string {
	return fixip.IpAddress
}

func (fixip *SFixedIP) GetIP() string {
	return fixip.IpAddress
}

func (fixip *SFixedIP) GetINetworkId() string {
	return fixip.NetworkId
}

func (fixip *SFixedIP) IsPrimary() bool {
	return true
}

type Port struct {
	multicloud.SNetworkInterfaceBase
	region          *SRegion
	ID              string `json:"id"`
	Name            string `json:"name"`
	Status          string `json:"status"`
	AdminStateUp    string `json:"admin_state_up"`
	DNSName         string `json:"dns_name"`
	MACAddress      string `json:"mac_address"`
	NetworkID       string `json:"network_id"`
	TenantID        string `json:"tenant_id"`
	DeviceID        string `json:"device_id"`
	DeviceOwner     string `json:"device_owner"`
	BindingVnicType string `json:"binding:vnic_type"`
	FixedIps        []SFixedIP
}

func (port *Port) GetName() string {
	if len(port.Name) > 0 {
		return port.Name
	}
	return port.ID
}

func (port *Port) GetId() string {
	return port.ID
}

func (port *Port) GetGlobalId() string {
	return port.ID
}

func (port *Port) GetMacAddress() string {
	return port.MACAddress
}

func (port *Port) GetAssociateType() string {
	switch port.DeviceOwner {
	case "compute:nova":
		return api.NETWORK_INTERFACE_ASSOCIATE_TYPE_SERVER
	case "network:router_gateway", "network:router_interface", "network:router_interface_distributed":
		return api.NETWORK_INTERFACE_ASSOCIATE_TYPE_RESERVED
	case "network:dhcp":
		return api.NETWORK_INTERFACE_ASSOCIATE_TYPE_DHCP
	case "neutron:LOADBALANCERV2":
		return api.NETWORK_INTERFACE_ASSOCIATE_TYPE_LOADBALANCER
	case "neutron:VIP_PORT":
		return api.NETWORK_INTERFACE_ASSOCIATE_TYPE_VIP
	}
	return port.DeviceOwner
}

func (port *Port) GetAssociateId() string {
	return port.DeviceID
}

func (port *Port) GetStatus() string {
	switch port.Status {
	case "ACTIVE", "DOWN":
		return api.NETWORK_INTERFACE_STATUS_AVAILABLE
	case "BUILD":
		return api.NETWORK_INTERFACE_STATUS_CREATING
	}
	return port.Status
}

func (port *Port) GetICloudInterfaceAddresses() ([]cloudprovider.ICloudInterfaceAddress, error) {
	address := []cloudprovider.ICloudInterfaceAddress{}
	for i := 0; i < len(port.FixedIps); i++ {
		port.FixedIps[i].NetworkId = port.NetworkID
		address = append(address, &port.FixedIps[i])
	}
	return address, nil
}

func (region *SRegion) GetINetworkInterfaces() ([]cloudprovider.ICloudNetworkInterface, error) {
	ports, err := region.GetPorts("")
	if err != nil {
		return nil, err
	}
	ret := []cloudprovider.ICloudNetworkInterface{}
	for i := 0; i < len(ports); i++ {
		if len(ports[i].DeviceID) == 0 || !utils.IsInStringArray(ports[i].DeviceOwner, []string{"compute:CCI", "compute:nova", "neutron:LOADBALANCERV2"}) {
			ports[i].region = region
			ret = append(ret, &ports[i])
		}
	}
	return ret, nil
}

func (self *SRegion) GetPort(portId string) (Port, error) {
	port := Port{}
	err := DoGet(self.ecsClient.Port.Get, portId, nil, &port)
	return port, err
}

// https://support.huaweicloud.com/api-vpc/zh-cn_topic_0030591299.html
func (self *SRegion) GetPorts(instanceId string) ([]Port, error) {
	ports := make([]Port, 0)
	querys := map[string]string{}
	if len(instanceId) > 0 {
		querys["device_id"] = instanceId
	}

	err := doListAllWithMarker(self.ecsClient.Port.List, querys, &ports)
	return ports, err
}

func (self *SEipAddress) GetProjectId() string {
	return ""
}

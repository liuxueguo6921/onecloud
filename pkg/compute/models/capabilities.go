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

package models

import (
	"context"
	"database/sql"
	"fmt"

	"yunion.io/x/jsonutils"
	"yunion.io/x/pkg/tristate"
	"yunion.io/x/pkg/utils"
	"yunion.io/x/sqlchemy"

	api "yunion.io/x/onecloud/pkg/apis/compute"
	"yunion.io/x/onecloud/pkg/cloudcommon/consts"
	"yunion.io/x/onecloud/pkg/cloudcommon/db"
	"yunion.io/x/onecloud/pkg/cloudcommon/policy"
	"yunion.io/x/onecloud/pkg/httperrors"
	"yunion.io/x/onecloud/pkg/mcclient"
	"yunion.io/x/onecloud/pkg/util/rbacutils"
)

type SCapabilities struct {
	Hypervisors        []string `json:",allowempty"`
	Brands             []string `json:",allowempty"`
	ResourceTypes      []string `json:",allowempty"`
	StorageTypes       []string `json:",allowempty"`
	DataStorageTypes   []string `json:",allowempty"`
	GPUModels          []string `json:",allowempty"`
	MinNicCount        int
	MaxNicCount        int
	MinDataDiskCount   int
	MaxDataDiskCount   int
	SchedPolicySupport bool
	Usable             bool
	PublicNetworkCount int
	Specs              jsonutils.JSONObject
}

func GetCapabilities(ctx context.Context, userCred mcclient.TokenCredential, query jsonutils.JSONObject, region *SCloudregion, zone *SZone) (SCapabilities, error) {
	capa := SCapabilities{}
	scopeStr := jsonutils.GetAnyString(query, []string{"scope"})
	scope := rbacutils.String2Scope(scopeStr)
	var domainId string
	domainStr := jsonutils.GetAnyString(query, []string{"domain", "domain_id", "project_domain", "project_domain_id"})
	if len(domainStr) > 0 {
		domain, err := db.TenantCacheManager.FetchDomainByIdOrName(ctx, domainStr)
		if err != nil {
			if err == sql.ErrNoRows {
				return capa, httperrors.NewResourceNotFoundError2("domains", domainStr)
			}
			return capa, httperrors.NewGeneralError(err)
		}
		domainId = domain.GetId()
	} else {
		domainId = userCred.GetProjectDomainId()
	}
	if scope == rbacutils.ScopeSystem {
		result := policy.PolicyManager.Allow(scope, userCred, consts.GetServiceType(), "capabilities", policy.PolicyActionList)
		if result != rbacutils.Allow {
			return capa, httperrors.NewForbiddenError("not allow to query system capability")
		}
		domainId = ""
	}
	capa.Hypervisors = getHypervisors(region, zone, domainId)
	capa.Brands = getBrands(region, zone, domainId, capa.Hypervisors)
	capa.ResourceTypes = getResourceTypes(region, zone, domainId)
	capa.StorageTypes = getStorageTypes(region, zone, true, domainId)
	capa.DataStorageTypes = getStorageTypes(region, zone, false, domainId)
	capa.GPUModels = getGPUs(region, zone, domainId)
	capa.SchedPolicySupport = isSchedPolicySupported(region, zone)
	capa.MinNicCount = getMinNicCount(region, zone)
	capa.MaxNicCount = getMaxNicCount(region, zone)
	capa.MinDataDiskCount = getMinDataDiskCount(region, zone)
	capa.MaxDataDiskCount = getMaxDataDiskCount(region, zone)
	capa.Usable = isUsable(region, zone, domainId)
	if query == nil {
		query = jsonutils.NewDict()
	}
	var err error
	if region != nil {
		query.(*jsonutils.JSONDict).Add(jsonutils.NewString(region.GetId()), "region")
	}
	if zone != nil {
		query.(*jsonutils.JSONDict).Add(jsonutils.NewString(zone.GetId()), "zone")
	}
	if len(domainId) > 0 {
		query.(*jsonutils.JSONDict).Add(jsonutils.NewString(domainId), "domain_id")
	}
	publicNetworkCount, _ := getNetworkPublicCount(region, zone, domainId)
	capa.PublicNetworkCount = publicNetworkCount
	mans := []ISpecModelManager{HostManager, IsolatedDeviceManager}
	capa.Specs, err = GetModelsSpecs(ctx, userCred, query.(*jsonutils.JSONDict), mans...)
	return capa, err
}

func getRegionZoneSubq(region *SCloudregion) *sqlchemy.SSubQuery {
	return ZoneManager.Query("id").Equals("cloudregion_id", region.GetId()).SubQuery()
}

func getDomainManagerSubq(domainId string) *sqlchemy.SSubQuery {
	providers := CloudproviderManager.Query().SubQuery()
	accounts := CloudaccountManager.Query().SubQuery()

	q := providers.Query(providers.Field("id"))
	q = q.Join(accounts, sqlchemy.Equals(accounts.Field("id"), providers.Field("cloudaccount_id")))
	q = q.Filter(sqlchemy.OR(
		sqlchemy.Equals(accounts.Field("domain_id"), domainId),
		sqlchemy.IsTrue(accounts.Field("is_public")),
	))
	q = q.Filter(sqlchemy.Equals(accounts.Field("status"), api.CLOUD_PROVIDER_CONNECTED))
	q = q.Filter(sqlchemy.IsTrue(accounts.Field("enabled")))

	return q.SubQuery()
}

func getBrands(region *SCloudregion, zone *SZone, domainId string, hypervisors []string) []string {
	q := CloudaccountManager.Query("brand").IsTrue("enabled")
	if zone != nil {
		region = zone.GetRegion()
	}
	if region != nil {
		providers := CloudproviderManager.Query().SubQuery()
		providerregions := CloudproviderRegionManager.Query().SubQuery()
		q = q.Join(providers, sqlchemy.Equals(q.Field("id"), providers.Field("cloudaccount_id")))
		q = q.Join(providerregions, sqlchemy.Equals(providers.Field("id"), providerregions.Field("cloudprovider_id")))
		q = q.Filter(sqlchemy.Equals(providerregions.Field("cloudregion_id"), region.Id))
	}
	if len(domainId) > 0 {
		q = q.Filter(sqlchemy.OR(
			sqlchemy.IsTrue(q.Field("is_public")),
			sqlchemy.Equals(q.Field("domain_id"), domainId),
		))
	}
	q = q.Distinct()
	rows, err := q.Rows()
	if err != nil {
		return nil
	}
	defer rows.Close()
	brands := make([]string, 0)
	for rows.Next() {
		var brand string
		rows.Scan(&brand)
		if len(brand) > 0 {
			brands = append(brands, brand)
		}
	}
	for _, hyper := range api.ONECLOUD_HYPERVISORS {
		if utils.IsInStringArray(hyper, hypervisors) {
			brands = append(brands, api.CLOUD_PROVIDER_ONECLOUD)
			break
		}
	}
	return brands
}

func getHypervisors(region *SCloudregion, zone *SZone, domainId string) []string {
	q := HostManager.Query("host_type", "manager_id")
	if region != nil {
		subq := getRegionZoneSubq(region)
		q = q.Filter(sqlchemy.In(q.Field("zone_id"), subq))
	}
	if zone != nil {
		q = q.Equals("zone_id", zone.Id)
	}
	if len(domainId) > 0 {
		subq := getDomainManagerSubq(domainId)
		q = q.Filter(sqlchemy.OR(
			sqlchemy.In(q.Field("manager_id"), subq),
			sqlchemy.IsNullOrEmpty(q.Field("manager_id")),
		))
	}
	q = q.IsNotEmpty("host_type").IsNotNull("host_type")
	// q = q.Equals("host_status", HOST_ONLINE)
	q = q.IsTrue("enabled")
	q = q.Distinct()
	rows, err := q.Rows()
	if err != nil {
		return nil
	}
	defer rows.Close()
	hypervisors := make([]string, 0)
	for rows.Next() {
		var hostType string
		var managerId string
		rows.Scan(&hostType, &managerId)
		if len(hostType) > 0 && IsProviderAccountEnabled(managerId) {
			hypervisor := api.HOSTTYPE_HYPERVISOR[hostType]
			if !utils.IsInStringArray(hypervisor, hypervisors) {
				hypervisors = append(hypervisors, hypervisor)
			}
		}
	}
	return hypervisors
}

func getResourceTypes(region *SCloudregion, zone *SZone, domainId string) []string {
	q := HostManager.Query("resource_type", "manager_id")
	if region != nil {
		subq := getRegionZoneSubq(region)
		q = q.Filter(sqlchemy.In(q.Field("zone_id"), subq))
	}
	if zone != nil {
		q = q.Equals("zone_id", zone.Id)
	}
	if len(domainId) > 0 {
		subq := getDomainManagerSubq(domainId)
		q = q.Filter(sqlchemy.OR(
			sqlchemy.In(q.Field("manager_id"), subq),
			sqlchemy.IsNullOrEmpty(q.Field("manager_id")),
		))
	}
	q = q.IsNotEmpty("resource_type").IsNotNull("resource_type")
	q = q.IsTrue("enabled")
	q = q.Distinct()
	rows, err := q.Rows()
	if err != nil {
		return nil
	}
	defer rows.Close()
	resourceTypes := make([]string, 0)
	for rows.Next() {
		var resType string
		var managerId string
		rows.Scan(&resType, &managerId)
		if len(resType) > 0 && IsProviderAccountEnabled(managerId) {
			if !utils.IsInStringArray(resType, resourceTypes) {
				resourceTypes = append(resourceTypes, resType)
			}
		}
	}
	return resourceTypes
}

func getStorageTypes(region *SCloudregion, zone *SZone, isSysDisk bool, domainId string) []string {
	storages := StorageManager.Query().SubQuery()
	hostStorages := HoststorageManager.Query().SubQuery()
	hosts := HostManager.Query().SubQuery()

	q := storages.Query(storages.Field("storage_type"), storages.Field("medium_type"))
	q = q.Join(hostStorages, sqlchemy.Equals(
		hostStorages.Field("storage_id"),
		storages.Field("id"),
	))
	q = q.Join(hosts, sqlchemy.Equals(
		hosts.Field("id"),
		hostStorages.Field("host_id"),
	))
	if region != nil {
		subq := getRegionZoneSubq(region)
		q = q.Filter(sqlchemy.In(storages.Field("zone_id"), subq))
	}
	if zone != nil {
		q = q.Filter(sqlchemy.Equals(storages.Field("zone_id"), zone.Id))
	}
	if len(domainId) > 0 {
		subq := getDomainManagerSubq(domainId)
		q = q.Filter(sqlchemy.OR(
			sqlchemy.In(hosts.Field("manager_id"), subq),
			sqlchemy.IsNullOrEmpty(hosts.Field("manager_id")),
		))
	}
	q = q.Filter(sqlchemy.Equals(hosts.Field("resource_type"), api.HostResourceTypeShared))
	q = q.Filter(sqlchemy.IsNotEmpty(storages.Field("storage_type")))
	q = q.Filter(sqlchemy.IsNotNull(storages.Field("storage_type")))
	q = q.Filter(sqlchemy.IsNotEmpty(storages.Field("medium_type")))
	q = q.Filter(sqlchemy.IsNotNull(storages.Field("medium_type")))
	q = q.Filter(sqlchemy.In(storages.Field("status"), []string{api.STORAGE_ENABLED, api.STORAGE_ONLINE}))
	q = q.Filter(sqlchemy.IsTrue(storages.Field("enabled")))
	if isSysDisk {
		q = q.Filter(sqlchemy.IsTrue(storages.Field("is_sys_disk_store")))
	}
	q = q.Filter(sqlchemy.NotEquals(hosts.Field("host_type"), api.HOST_TYPE_BAREMETAL))
	q = q.Distinct()
	rows, err := q.Rows()
	if err != nil {
		return nil
	}
	defer rows.Close()
	storageTypes := make([]string, 0)
	for rows.Next() {
		var storageType, mediumType string
		rows.Scan(&storageType, &mediumType)
		if len(storageType) > 0 && len(mediumType) > 0 {
			storageTypes = append(storageTypes, fmt.Sprintf("%s/%s", storageType, mediumType))
		}
	}
	return storageTypes
}

func getGPUs(region *SCloudregion, zone *SZone, domainId string) []string {
	devices := IsolatedDeviceManager.Query().SubQuery()
	hosts := HostManager.Query().SubQuery()

	q := devices.Query(devices.Field("model"))
	if region != nil {
		subq := getRegionZoneSubq(region)
		q = q.Join(hosts, sqlchemy.Equals(devices.Field("host_id"), hosts.Field("id")))
		q = q.Filter(sqlchemy.In(hosts.Field("zone_id"), subq))
	}
	if zone != nil {
		q = q.Join(hosts, sqlchemy.Equals(devices.Field("host_id"), hosts.Field("id")))
		q = q.Filter(sqlchemy.Equals(hosts.Field("zone_id"), zone.Id))
	}
	if len(domainId) > 0 {
		subq := getDomainManagerSubq(domainId)
		q = q.Filter(sqlchemy.OR(
			sqlchemy.In(hosts.Field("manager_id"), subq),
			sqlchemy.IsNullOrEmpty(hosts.Field("manager_id")),
		))
	}
	q = q.Distinct()

	rows, err := q.Rows()
	if err != nil {
		return nil
	}
	defer rows.Close()
	gpus := make([]string, 0)
	for rows.Next() {
		var model string
		rows.Scan(&model)
		if len(model) > 0 {
			gpus = append(gpus, model)
		}
	}
	return gpus
}

func getNetworkCount(region *SCloudregion, zone *SZone, domainId string) (int, error) {
	return getNetworkCountByFilter(region, zone, domainId, tristate.None)
}

func getNetworkPublicCount(region *SCloudregion, zone *SZone, domainId string) (int, error) {
	return getNetworkCountByFilter(region, zone, domainId, tristate.True)
}

func getNetworkCountByFilter(region *SCloudregion, zone *SZone, domainId string, isPublic tristate.TriState) (int, error) {
	if zone != nil && region == nil {
		region = zone.GetRegion()
	}

	vpcs := VpcManager.Query().SubQuery()
	wires := WireManager.Query().SubQuery()
	networks := NetworkManager.Query().SubQuery()

	q := networks.Query()
	if !isPublic.IsNone() {
		if isPublic.IsTrue() {
			q = q.IsTrue("is_public")
		} else {
			q = q.IsFalse("is_public")
		}
	}
	q = q.Join(wires, sqlchemy.Equals(networks.Field("wire_id"), wires.Field("id")))

	if region != nil {
		if utils.IsInStringArray(region.Provider, api.REGIONAL_NETWORK_PROVIDERS) {
			wires := WireManager.Query().SubQuery()
			vpcs := VpcManager.Query().SubQuery()
			subq := wires.Query(wires.Field("id")).
				Join(vpcs, sqlchemy.Equals(wires.Field("vpc_id"), vpcs.Field("id"))).
				Filter(sqlchemy.Equals(vpcs.Field("cloudregion_id"), region.Id))
			q = q.Filter(sqlchemy.In(q.Field("wire_id"), subq))
		} else {
			subq := getRegionZoneSubq(region)
			q = q.Filter(sqlchemy.In(wires.Field("zone_id"), subq))
		}
	}
	if zone != nil && !utils.IsInStringArray(region.Provider, api.REGIONAL_NETWORK_PROVIDERS) {
		q = q.Filter(sqlchemy.Equals(wires.Field("zone_id"), zone.Id))
	}
	if len(domainId) > 0 {
		subq := getDomainManagerSubq(domainId)
		q = q.Join(vpcs, sqlchemy.Equals(wires.Field("vpc_id"), vpcs.Field("id")))
		q = q.Filter(sqlchemy.OR(
			sqlchemy.In(vpcs.Field("manager_id"), subq),
			sqlchemy.IsNullOrEmpty(vpcs.Field("manager_id")),
		))
		if isPublic.Bool() {
			q = q.Filter(sqlchemy.OR(
				sqlchemy.Equals(q.Field("public_scope"), rbacutils.ScopeSystem),
				sqlchemy.AND(
					sqlchemy.Equals(q.Field("public_scope"), rbacutils.ScopeDomain),
					sqlchemy.Equals(q.Field("domain_id"), domainId))))
		}
	}
	q = q.Filter(sqlchemy.Equals(networks.Field("status"), api.NETWORK_STATUS_AVAILABLE))

	return q.CountWithError()
}

func isSchedPolicySupported(region *SCloudregion, zone *SZone) bool {
	return true
}

func getMinNicCount(region *SCloudregion, zone *SZone) int {
	if region != nil {
		return region.getMinNicCount()
	}
	if zone != nil {
		return zone.getMinNicCount()
	}
	return 0
}

func getMaxNicCount(region *SCloudregion, zone *SZone) int {
	if region != nil {
		return region.getMaxNicCount()
	}
	if zone != nil {
		return zone.getMaxNicCount()
	}
	return 0
}

func getMinDataDiskCount(region *SCloudregion, zone *SZone) int {
	if region != nil {
		return region.getMinDataDiskCount()
	}
	if zone != nil {
		return zone.getMinDataDiskCount()
	}
	return 0
}

func getMaxDataDiskCount(region *SCloudregion, zone *SZone) int {
	if region != nil {
		return region.getMaxDataDiskCount()
	}
	if zone != nil {
		return zone.getMaxDataDiskCount()
	}
	return 0
}

func isUsable(region *SCloudregion, zone *SZone, domainId string) bool {
	cnt, err := getNetworkCount(region, zone, domainId)
	if err != nil {
		return false
	}
	if cnt > 0 {
		return true
	} else {
		return false
	}
}

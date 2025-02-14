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

package regiondrivers

import (
	"context"
	"fmt"

	"yunion.io/x/jsonutils"

	api "yunion.io/x/onecloud/pkg/apis/compute"
	"yunion.io/x/onecloud/pkg/cloudcommon/db/taskman"
	"yunion.io/x/onecloud/pkg/cloudprovider"
	"yunion.io/x/onecloud/pkg/compute/models"
	"yunion.io/x/onecloud/pkg/httperrors"
	"yunion.io/x/onecloud/pkg/mcclient"
)

type SBaseRegionDriver struct {
}

func (self *SBaseRegionDriver) RequestCreateLoadbalancer(ctx context.Context, userCred mcclient.TokenCredential, lb *models.SLoadbalancer, task taskman.ITask) error {
	return fmt.Errorf("Not Implement RequestCreateLoadbalancer")
}

func (self *SBaseRegionDriver) RequestStartLoadbalancer(ctx context.Context, userCred mcclient.TokenCredential, lb *models.SLoadbalancer, task taskman.ITask) error {
	return fmt.Errorf("Not Implement RequestStartLoadbalancer")
}

func (self *SBaseRegionDriver) RequestStopLoadbalancer(ctx context.Context, userCred mcclient.TokenCredential, lb *models.SLoadbalancer, task taskman.ITask) error {
	return fmt.Errorf("Not Implement RequestStopLoadbalancer")
}

func (self *SBaseRegionDriver) RequestSyncstatusLoadbalancer(ctx context.Context, userCred mcclient.TokenCredential, lb *models.SLoadbalancer, task taskman.ITask) error {
	return fmt.Errorf("Not Implement RequestSyncstatusLoadbalancer")
}

func (self *SBaseRegionDriver) RequestDeleteLoadbalancer(ctx context.Context, userCred mcclient.TokenCredential, lb *models.SLoadbalancer, task taskman.ITask) error {
	return fmt.Errorf("Not Implement RequestDeleteLoadbalancer")
}

func (self *SBaseRegionDriver) RequestCreateLoadbalancerAcl(ctx context.Context, userCred mcclient.TokenCredential, lbacl *models.SCachedLoadbalancerAcl, task taskman.ITask) error {
	return fmt.Errorf("Not Implement RequestCreateLoadbalancerAcl")
}

func (self *SBaseRegionDriver) RequestSyncLoadbalancerAcl(ctx context.Context, userCred mcclient.TokenCredential, lbacl *models.SCachedLoadbalancerAcl, task taskman.ITask) error {
	return fmt.Errorf("Not Implement RequestSyncLoadbalancerAcl")
}

func (self *SBaseRegionDriver) RequestDeleteLoadbalancerAcl(ctx context.Context, userCred mcclient.TokenCredential, lbacl *models.SCachedLoadbalancerAcl, task taskman.ITask) error {
	return fmt.Errorf("Not Implement RequestDeleteLoadbalancerAcl")
}

func (self *SBaseRegionDriver) RequestCreateLoadbalancerCertificate(ctx context.Context, userCred mcclient.TokenCredential, lbcert *models.SCachedLoadbalancerCertificate, task taskman.ITask) error {
	return fmt.Errorf("Not Implement RequestCreateLoadbalancerCertificate")
}

func (self *SBaseRegionDriver) RequestDeleteLoadbalancerCertificate(ctx context.Context, userCred mcclient.TokenCredential, lbcert *models.SCachedLoadbalancerCertificate, task taskman.ITask) error {
	return fmt.Errorf("Not Implement RequestDeleteLoadbalancerCertificate")
}

func (self *SBaseRegionDriver) RequestCreateLoadbalancerBackendGroup(ctx context.Context, userCred mcclient.TokenCredential, lbbg *models.SLoadbalancerBackendGroup, backends []cloudprovider.SLoadbalancerBackend, task taskman.ITask) error {
	return fmt.Errorf("Not Implement RequestCreateLoadbalancerBackendGroup")
}

func (self *SBaseRegionDriver) RequestDeleteLoadbalancerBackendGroup(ctx context.Context, userCred mcclient.TokenCredential, lbbg *models.SLoadbalancerBackendGroup, task taskman.ITask) error {
	return fmt.Errorf("Not Implement RequestDeleteLoadbalancerBackendGroup")
}

func (self *SBaseRegionDriver) RequestSyncLoadbalancerBackendGroup(ctx context.Context, userCred mcclient.TokenCredential, lblis *models.SLoadbalancerListener, lbbg *models.SLoadbalancerBackendGroup, task taskman.ITask) error {
	return fmt.Errorf("Not Implement RequestSyncLoadbalancerBackendGroup")
}

func (self *SBaseRegionDriver) RequestCreateLoadbalancerBackend(ctx context.Context, userCred mcclient.TokenCredential, lbb *models.SLoadbalancerBackend, task taskman.ITask) error {
	return fmt.Errorf("Not Implement RequestCreateLoadbalancerBackend")
}

func (self *SBaseRegionDriver) RequestDeleteLoadbalancerBackend(ctx context.Context, userCred mcclient.TokenCredential, lbb *models.SLoadbalancerBackend, task taskman.ITask) error {
	return fmt.Errorf("Not Implement RequestDeleteLoadbalancerBackend")
}

func (self *SBaseRegionDriver) RequestSyncLoadbalancerBackend(ctx context.Context, userCred mcclient.TokenCredential, lbb *models.SLoadbalancerBackend, task taskman.ITask) error {
	return fmt.Errorf("Not Implement RequestSyncLoadbalancerBackend")
}

func (self *SBaseRegionDriver) RequestCreateLoadbalancerListener(ctx context.Context, userCred mcclient.TokenCredential, lblis *models.SLoadbalancerListener, task taskman.ITask) error {
	return fmt.Errorf("Not Implement RequestCreateLoadbalancerListener")
}

func (self *SBaseRegionDriver) RequestDeleteLoadbalancerListener(ctx context.Context, userCred mcclient.TokenCredential, lblis *models.SLoadbalancerListener, task taskman.ITask) error {
	return fmt.Errorf("Not Implement RequestDeleteLoadbalancerListener")
}

func (self *SBaseRegionDriver) RequestStartLoadbalancerListener(ctx context.Context, userCred mcclient.TokenCredential, lblis *models.SLoadbalancerListener, task taskman.ITask) error {
	return fmt.Errorf("Not Implement RequestStartLoadbalancerListener")
}

func (self *SBaseRegionDriver) RequestStopLoadbalancerListener(ctx context.Context, userCred mcclient.TokenCredential, lblis *models.SLoadbalancerListener, task taskman.ITask) error {
	return fmt.Errorf("Not Implement RequestStopLoadbalancerListener")
}

func (self *SBaseRegionDriver) RequestSyncstatusLoadbalancerListener(ctx context.Context, userCred mcclient.TokenCredential, lblis *models.SLoadbalancerListener, task taskman.ITask) error {
	return fmt.Errorf("Not Implement RequestSyncstatusLoadbalancerListener")
}

func (self *SBaseRegionDriver) RequestSyncLoadbalancerListener(ctx context.Context, userCred mcclient.TokenCredential, lblis *models.SLoadbalancerListener, task taskman.ITask) error {
	return fmt.Errorf("Not Implement RequestSyncLoadbalancerListener")
}

func (self *SBaseRegionDriver) RequestCreateLoadbalancerListenerRule(ctx context.Context, userCred mcclient.TokenCredential, lbr *models.SLoadbalancerListenerRule, task taskman.ITask) error {
	return fmt.Errorf("Not Implement RequestCreateLoadbalancerListenerRule")
}

func (self *SBaseRegionDriver) RequestDeleteLoadbalancerListenerRule(ctx context.Context, userCred mcclient.TokenCredential, lbr *models.SLoadbalancerListenerRule, task taskman.ITask) error {
	return fmt.Errorf("Not Implement RequestDeleteLoadbalancerListenerRule")
}

func (self *SBaseRegionDriver) ValidateCreateSnapshotPolicyData(ctx context.Context, userCred mcclient.TokenCredential, input *api.SSnapshotPolicyCreateInput, ownerId mcclient.IIdentityProvider, data *jsonutils.JSONDict) error {
	var err error

	if len(input.RepeatWeekdays) == 0 {
		return httperrors.NewMissingParameterError("repeat_weekdays")
	}
	input.RepeatWeekdays, err = daysValidate(input.RepeatWeekdays, 1, 7)
	if err != nil {
		return httperrors.NewInputParameterError(err.Error())
	}

	if len(input.TimePoints) == 0 {
		return httperrors.NewInputParameterError("time_points")
	}
	input.TimePoints, err = daysValidate(input.TimePoints, 0, 23)
	if err != nil {
		return httperrors.NewInputParameterError(err.Error())
	}
	return nil
}

func (self *SBaseRegionDriver) RequestCreateSnapshotPolicy(ctx context.Context, userCred mcclient.TokenCredential, sp *models.SSnapshotPolicy, task taskman.ITask) error {
	return fmt.Errorf("Not Implement RequestCreateSnapshotPolicy")
}

func (self *SBaseRegionDriver) RequestDeleteSnapshotPolicy(ctx context.Context, userCred mcclient.TokenCredential, sp *models.SSnapshotPolicy, task taskman.ITask) error {
	return fmt.Errorf("Not Implement RequestDeleteSnapshotPolicy")
}

func (self *SBaseRegionDriver) RequestApplySnapshotPolicy(ctx context.Context, userCred mcclient.TokenCredential, sp *models.SSnapshotPolicy, task taskman.ITask, diskId string) error {
	return fmt.Errorf("Not Implement RequestApplySnapshotPolicy")
}

func (self *SBaseRegionDriver) RequestCancelSnapshotPolicy(ctx context.Context, userCred mcclient.TokenCredential, sp *models.SSnapshotPolicy, task taskman.ITask, diskId string) error {
	return fmt.Errorf("Not Implement RequestApplySnapshotPolicy")
}

func (self *SBaseRegionDriver) ValidateSnapshotDelete(ctx context.Context, snapshot *models.SSnapshot) error {
	return fmt.Errorf("Not Implement ValidateSnapshotDelete")
}

func (self *SBaseRegionDriver) RequestDeleteSnapshot(ctx context.Context, snapshot *models.SSnapshot, task taskman.ITask) error {
	return fmt.Errorf("Not Implement RequestDeleteSnapshot")
}

func (self *SBaseRegionDriver) ValidateSnapshotCreate(ctx context.Context, userCred mcclient.TokenCredential, disk *models.SDisk, data *jsonutils.JSONDict) error {
	return fmt.Errorf("Not Implement ValidateSnapshotCreate")
}

func (self *SBaseRegionDriver) RequestCreateSnapshot(ctx context.Context, snapshot *models.SSnapshot, task taskman.ITask) error {
	return fmt.Errorf("Not Implement RequestCreateSnapshot")
}

func (self *SBaseRegionDriver) SnapshotIsOutOfChain(disk *models.SDisk) bool {
	return true
}

func (self *SBaseRegionDriver) GetDiskResetParams(snapshot *models.SSnapshot) *jsonutils.JSONDict {
	return nil
}

func (self *SBaseRegionDriver) OnDiskReset(ctx context.Context, userCred mcclient.TokenCredential, disk *models.SDisk, snapshot *models.SSnapshot, data *jsonutils.JSONObject) error {
	return fmt.Errorf("Not Implement OnDiskReset")
}

func (self *SBaseRegionDriver) ValidateCreateSnapshopolicyDiskData(ctx context.Context, userCred mcclient.TokenCredential, diskID string) error {
	return fmt.Errorf("Not Implement ValidateCreateSnapshotpolicyDiskData")
}

func (self *SBaseRegionDriver) OnSnapshotDelete(ctx context.Context, snapshot *models.SSnapshot, task taskman.ITask, data jsonutils.JSONObject) error {
	return fmt.Errorf("Not implement OnSnapshotDelete")
}

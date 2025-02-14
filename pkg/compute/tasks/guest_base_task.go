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

package tasks

import (
	"context"

	"yunion.io/x/onecloud/pkg/cloudcommon/db/taskman"
	"yunion.io/x/onecloud/pkg/compute/models"
	"yunion.io/x/onecloud/pkg/util/rbacutils"
)

type SGuestBaseTask struct {
	taskman.STask
}

func (self *SGuestBaseTask) getGuest() *models.SGuest {
	obj := self.GetObject()
	return obj.(*models.SGuest)
}

func (self *SGuestBaseTask) SetStageFailed(ctx context.Context, reason string) {
	self.finalReleasePendingUsage(ctx)
	self.STask.SetStageFailed(ctx, reason)
}

func (self *SGuestBaseTask) finalReleasePendingUsage(ctx context.Context) {
	pendingUsage := models.SQuota{}
	err := self.GetPendingUsage(&pendingUsage)
	if err == nil && !pendingUsage.IsEmpty() {
		guest := self.getGuest()
		quotaPlatform := guest.GetQuotaPlatformID()
		models.QuotaManager.CancelPendingUsage(ctx, self.UserCred, rbacutils.ScopeProject, guest.GetOwnerId(), quotaPlatform, &pendingUsage, &pendingUsage)
	}
}

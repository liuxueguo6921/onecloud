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

package modules

import (
	"fmt"

	"yunion.io/x/onecloud/pkg/mcclient"
)

type GroupManager struct {
	ResourceManager
}

func (this *GroupManager) GetUsers(s *mcclient.ClientSession, gid string) (*ListResult, error) {
	url := fmt.Sprintf("/groups/%s/users", gid)
	return this._list(s, url, "users")
}

var (
	Groups GroupManager
)

func (this *GroupManager) GetProjects(session *mcclient.ClientSession, uid string) (*ListResult, error) {
	url := fmt.Sprintf("/groups/%s/projects?admin=true", uid)
	return this._list(session, url, "projects")
}

func init() {
	Groups = GroupManager{NewIdentityV3Manager("group", "groups",
		[]string{},
		[]string{"ID", "Name", "Domain_Id", "project_domain",
			"User_Count", "Description"})}

	register(&Groups)
}

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

package options

import (
	"yunion.io/x/onecloud/pkg/cloudcommon/options"
)

type NotifyOption struct {
	options.CommonOptions
	options.DBOptions

	DingtalkEnabled bool   `help:"Enable dingtalk"`
	SocketFileDir   string `help:"Socket file directory" default:"/etc/yunion/notify"`
	UpdateInterval  int    `help:"Update send services interval(unit:s)" default:30`
	VerifyEmailUrl  string
	ReSendScope     int `help:"Resend all messages that have not been sent successfully within ReSendScope minutes"`
}

var Options NotifyOption

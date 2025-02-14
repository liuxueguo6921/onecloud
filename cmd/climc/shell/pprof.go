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

package shell

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"syscall"

	"yunion.io/x/pkg/util/signalutils"

	"yunion.io/x/onecloud/pkg/mcclient"
	"yunion.io/x/onecloud/pkg/mcclient/modules"
	"yunion.io/x/onecloud/pkg/util/netutils2"
	"yunion.io/x/onecloud/pkg/util/procutils"
)

func init() {
	type TraceOptions struct {
		Second  int    `help:"pprof seconds" short-token:"s"`
		SERVICE string `help:"Service type"`
	}

	downloadToTemp := func(input io.Reader, pattern string) (string, error) {
		tmpfile, err := ioutil.TempFile("", pattern)
		if err != nil {
			return "", err
		}
		defer tmpfile.Close()
		if _, err := io.Copy(tmpfile, input); err != nil {
			return "", err
		}
		return tmpfile.Name(), nil
	}

	pprofRun := func(s *mcclient.ClientSession, svcType, pType string, second int, args ...string) error {
		src, err := modules.GetPProfByType(s, svcType, pType, second)
		if err != nil {
			return err
		}
		tempfile, err := downloadToTemp(src, pType)
		if err != nil {
			return err
		}
		defer func() { os.Remove(tempfile) }()

		signalutils.RegisterSignal(func() {
			os.Remove(tempfile)
			os.Exit(0)
		}, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
		signalutils.StartTrap()

		cmd := procutils.NewCommand("go", "tool")
		cmd.Args = append(cmd.Args, args...)
		cmd.Args = append(cmd.Args, tempfile)
		if _, err := cmd.Run(); err != nil {
			return err
		}
		return nil
	}

	R(&TraceOptions{}, "pprof-trace", "pprof trace of backend service", func(s *mcclient.ClientSession, args *TraceOptions) error {
		return pprofRun(s, args.SERVICE, "trace", args.Second, "trace")
	})

	R(&TraceOptions{}, "pprof-profile", "pprof profile of backend service", func(s *mcclient.ClientSession, args *TraceOptions) error {
		port, err := netutils2.GetFreePort()
		if err != nil {
			return err
		}
		return pprofRun(s, args.SERVICE, "profile", args.Second, "pprof", fmt.Sprintf("-http=:%d", port))
	})
}

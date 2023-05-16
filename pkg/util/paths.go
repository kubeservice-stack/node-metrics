// Copyright 2015 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"path/filepath"
	"strings"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/procfs"
)

var (
	// The path of the proc filesystem.
	ProcPath     = kingpin.Flag("path.procfs", "procfs mountpoint.").Default(procfs.DefaultMountPoint).String()
	SysPath      = kingpin.Flag("path.sysfs", "sysfs mountpoint.").Default("/sys").String()
	RootfsPath   = kingpin.Flag("path.rootfs", "rootfs mountpoint.").Default("/").String()
	UdevDataPath = kingpin.Flag("path.udev.data", "udev data path.").Default("/run/udev/data").String()
)

func ProcFilePath(name string) string {
	return filepath.Join(*ProcPath, name)
}

func SysFilePath(name string) string {
	return filepath.Join(*SysPath, name)
}

func RootfsFilePath(name string) string {
	return filepath.Join(*RootfsPath, name)
}

func UdevDataFilePath(name string) string {
	return filepath.Join(*UdevDataPath, name)
}

func RootfsStripPrefix(path string) string {
	if *RootfsPath == "/" {
		return path
	}
	stripped := strings.TrimPrefix(path, *RootfsPath)
	if stripped == "" {
		return "/"
	}
	return stripped
}

/*
Copyright © 2021 SUSE LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import (
	"errors"
	"fmt"
	"github.com/rancher-sandbox/elemental/pkg/types/v1"
	"github.com/spf13/afero"
	"github.com/zloylos/grsync"
	"io"
	"os/exec"
	"strings"
	"time"
)

func CommandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// BootedFrom will check if we are booting from the given label
func BootedFrom(runner v1.Runner, label string) bool {
	out, _ := runner.Run("cat", "/proc/cmdline")
	return strings.Contains(string(out), label)
}

// GetDeviceByLabel will try to return the device that matches the given label.
// attempts value sets the number of attempts to find the device, it
// waits a second between attempts.
func GetDeviceByLabel(runner v1.Runner, label string, attempts int) (string, error) {
	for tries := 0; tries < attempts; tries++ {
		runner.Run("udevadm", "settle")
		out, err := runner.Run("blkid", "--label", label)
		if err == nil && strings.TrimSpace(string(out)) != "" {
			return strings.TrimSpace(string(out)), nil
		}
		time.Sleep(1 * time.Second)
	}
	return "", errors.New("no device found")
}

// Copies source file to target file using afero.Fs interface
func CopyFile(fs afero.Fs, source string, target string) (err error) {
	sourceFile, err := fs.Open(source)
	if err != nil {
		return err
	}
	defer func() {
		if err == nil {
			err = sourceFile.Close()
		}
	}()

	targetFile, err := fs.Create(target)
	if err != nil {
		return err
	}
	defer func() {
		if err == nil {
			err = targetFile.Close()
		}
	}()

	_, err = io.Copy(targetFile, sourceFile)
	return err
}

// Copies source file to target file using afero.Fs interface
func CreateDirStructure(fs afero.Fs, target string) error {
	for _, dir := range []string{"sys", "proc", "dev", "tmp", "boot", "usr/local", "oem"} {
		err := fs.MkdirAll(fmt.Sprintf("%s/%s", target, dir), 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

// SyncData rsync's source folder contents to a target folder content,
// both are expected to exist before hand.
func SyncData(source string, target string, excludes ...string) error {
	if strings.HasSuffix(source, "/") == false {
		source = fmt.Sprintf("%s/", source)
	}

	if strings.HasSuffix(target, "/") == false {
		target = fmt.Sprintf("%s/", target)
	}

	task := grsync.NewTask(
		source,
		target,
		grsync.RsyncOptions{
			Quiet:   false,
			Archive: true,
			XAttrs:  true,
			ACLs:    true,
			Exclude: excludes,
		},
	)

	return task.Run()
}
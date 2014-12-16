/*
   Copyright 2014 CoreOS, Inc.

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

package configdrive

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

const (
	openstackApiVersion = "latest"
)

type configDrive struct {
	root     string
	readFile func(filename string) ([]byte, error)
}

func NewDatasource(root string) *configDrive {
	return &configDrive{root, ioutil.ReadFile}
}

func (cd *configDrive) IsAvailable() bool {
	_, err := os.Stat(cd.root)
	return !os.IsNotExist(err)
}

func (cd *configDrive) AvailabilityChanges() bool {
	return true
}

func (cd *configDrive) ConfigRoot() string {
	return cd.openstackRoot()
}

func (cd *configDrive) FetchMetadata() ([]byte, error) {
	return cd.tryReadFile(path.Join(cd.openstackVersionRoot(), "meta_data.json"))
}

func (cd *configDrive) FetchUserdata() ([]byte, error) {
	return cd.tryReadFile(path.Join(cd.openstackVersionRoot(), "user_data"))
}

func (cd *configDrive) FetchNetworkConfig(filename string) ([]byte, error) {
	if filename == "" {
		return []byte{}, nil
	}
	return cd.tryReadFile(path.Join(cd.openstackRoot(), filename))
}

func (cd *configDrive) Type() string {
	return "cloud-drive"
}

func (cd *configDrive) openstackRoot() string {
	return path.Join(cd.root, "lvmcloud")
}

func (cd *configDrive) openstackVersionRoot() string {
	return path.Join(cd.openstackRoot(), openstackApiVersion)
}

func (cd *configDrive) tryReadFile(filename string) ([]byte, error) {
	fmt.Printf("Attempting to read from %q\n", filename)
	data, err := cd.readFile(filename)
	if os.IsNotExist(err) {
		err = nil
	}
	return data, err
}

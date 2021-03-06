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

package validate

import (
	"reflect"
	"testing"
)

func TestCheckStructure(t *testing.T) {
	tests := []struct {
		config string

		entries []Entry
	}{
		{},

		// Test for unrecognized keys
		{
			config:  "test:",
			entries: []Entry{{entryWarning, "unrecognized key \"test\"", 1}},
		},
		{
			config:  "coreos:\n  etcd:\n    bad:",
			entries: []Entry{{entryWarning, "unrecognized key \"bad\"", 3}},
		},
		{
			config: "coreos:\n  etcd:\n    discovery: good",
		},

		// Test for error on list of nodes
		{
			config: "coreos:\n  units:\n    - hello\n    - goodbye",
			entries: []Entry{
				{entryWarning, "incorrect type for \"units[0]\" (want struct)", 3},
				{entryWarning, "incorrect type for \"units[1]\" (want struct)", 4},
			},
		},

		// Test for incorrect types
		// Want boolean
		{
			config: "coreos:\n  units:\n    - enable: true",
		},
		{
			config:  "coreos:\n  units:\n    - enable: 4",
			entries: []Entry{{entryWarning, "incorrect type for \"enable\" (want bool)", 3}},
		},
		{
			config:  "coreos:\n  units:\n    - enable: bad",
			entries: []Entry{{entryWarning, "incorrect type for \"enable\" (want bool)", 3}},
		},
		{
			config:  "coreos:\n  units:\n    - enable:\n        bad:",
			entries: []Entry{{entryWarning, "incorrect type for \"enable\" (want bool)", 3}},
		},
		{
			config:  "coreos:\n  units:\n    - enable:\n      - bad",
			entries: []Entry{{entryWarning, "incorrect type for \"enable\" (want bool)", 3}},
		},
		// Want string
		{
			config: "hostname: true",
		},
		{
			config: "hostname: 4",
		},
		{
			config: "hostname: host",
		},
		{
			config:  "hostname:\n  name:",
			entries: []Entry{{entryWarning, "incorrect type for \"hostname\" (want string)", 1}},
		},
		{
			config:  "hostname:\n  - name",
			entries: []Entry{{entryWarning, "incorrect type for \"hostname\" (want string)", 1}},
		},
		// Want struct
		{
			config:  "coreos: true",
			entries: []Entry{{entryWarning, "incorrect type for \"coreos\" (want struct)", 1}},
		},
		{
			config:  "coreos: 4",
			entries: []Entry{{entryWarning, "incorrect type for \"coreos\" (want struct)", 1}},
		},
		{
			config:  "coreos: hello",
			entries: []Entry{{entryWarning, "incorrect type for \"coreos\" (want struct)", 1}},
		},
		{
			config: "coreos:\n  etcd:\n    discovery: fire in the disco",
		},
		{
			config:  "coreos:\n  - hello",
			entries: []Entry{{entryWarning, "incorrect type for \"coreos\" (want struct)", 1}},
		},
		// Want []string
		{
			config:  "ssh_authorized_keys: true",
			entries: []Entry{{entryWarning, "incorrect type for \"ssh_authorized_keys\" (want []string)", 1}},
		},
		{
			config:  "ssh_authorized_keys: 4",
			entries: []Entry{{entryWarning, "incorrect type for \"ssh_authorized_keys\" (want []string)", 1}},
		},
		{
			config:  "ssh_authorized_keys: key",
			entries: []Entry{{entryWarning, "incorrect type for \"ssh_authorized_keys\" (want []string)", 1}},
		},
		{
			config:  "ssh_authorized_keys:\n  key: value",
			entries: []Entry{{entryWarning, "incorrect type for \"ssh_authorized_keys\" (want []string)", 1}},
		},
		{
			config: "ssh_authorized_keys:\n  - key",
		},
		{
			config:  "ssh_authorized_keys:\n  - key: value",
			entries: []Entry{{entryWarning, "incorrect type for \"ssh_authorized_keys[0]\" (want string)", 2}},
		},
		// Want []struct
		{
			config:  "users:\n  true",
			entries: []Entry{{entryWarning, "incorrect type for \"users\" (want []struct)", 1}},
		},
		{
			config:  "users:\n  4",
			entries: []Entry{{entryWarning, "incorrect type for \"users\" (want []struct)", 1}},
		},
		{
			config:  "users:\n  bad",
			entries: []Entry{{entryWarning, "incorrect type for \"users\" (want []struct)", 1}},
		},
		{
			config:  "users:\n  bad:",
			entries: []Entry{{entryWarning, "incorrect type for \"users\" (want []struct)", 1}},
		},
		{
			config: "users:\n  - name: good",
		},
		// Want struct within array
		{
			config:  "users:\n  - true",
			entries: []Entry{{entryWarning, "incorrect type for \"users[0]\" (want struct)", 2}},
		},
		{
			config:  "users:\n  - name: hi\n  - true",
			entries: []Entry{{entryWarning, "incorrect type for \"users[1]\" (want struct)", 3}},
		},
		{
			config:  "users:\n  - 4",
			entries: []Entry{{entryWarning, "incorrect type for \"users[0]\" (want struct)", 2}},
		},
		{
			config:  "users:\n  - bad",
			entries: []Entry{{entryWarning, "incorrect type for \"users[0]\" (want struct)", 2}},
		},
		{
			config:  "users:\n  - - bad",
			entries: []Entry{{entryWarning, "incorrect type for \"users[0]\" (want struct)", 2}},
		},
	}

	for i, tt := range tests {
		r := Report{}
		n, err := parseCloudConfig([]byte(tt.config), &r)
		if err != nil {
			panic(err)
		}
		checkStructure(n, &r)

		if e := r.Entries(); !reflect.DeepEqual(tt.entries, e) {
			t.Errorf("bad report (%d, %q): want %#v, got %#v", i, tt.config, tt.entries, e)
		}
	}
}

func TestCheckValidity(t *testing.T) {
	tests := []struct {
		config string

		entries []Entry
	}{
		// string
		{
			config: "hostname: test",
		},

		// int
		{
			config: "coreos:\n  fleet:\n    verbosity: 2",
		},

		// bool
		{
			config: "coreos:\n  units:\n    - enable: true",
		},

		// slice
		{
			config: "coreos:\n  units:\n    - command: start\n    - name: stop",
		},
		{
			config:  "coreos:\n  units:\n    - command: lol",
			entries: []Entry{{entryError, "invalid value lol", 3}},
		},

		// struct
		{
			config: "coreos:\n  update:\n    reboot_strategy: off",
		},
		{
			config:  "coreos:\n  update:\n    reboot_strategy: always",
			entries: []Entry{{entryError, "invalid value always", 3}},
		},

		// unknown
		{
			config: "unknown: hi",
		},
	}

	for i, tt := range tests {
		r := Report{}
		n, err := parseCloudConfig([]byte(tt.config), &r)
		if err != nil {
			panic(err)
		}
		checkValidity(n, &r)

		if e := r.Entries(); !reflect.DeepEqual(tt.entries, e) {
			t.Errorf("bad report (%d, %q): want %#v, got %#v", i, tt.config, tt.entries, e)
		}
	}
}

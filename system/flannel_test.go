package system

import (
	"reflect"
	"testing"

	"github.com/coreos/coreos-cloudinit/config"
)

func TestFlannelUnits(t *testing.T) {
	for _, tt := range []struct {
		config config.Flannel
		units  []Unit
	}{
		{
			config.Flannel{},
			[]Unit{{config.Unit{
				Name:    "flanneld.service",
				Runtime: true,
				DropIns: []config.UnitDropIn{{Name: "20-cloudinit.conf"}},
			}}},
		},
		{
			config.Flannel{
				EtcdEndpoint: "http://12.34.56.78:4001",
				EtcdPrefix:   "/coreos.com/network/tenant1",
			},
			[]Unit{{config.Unit{
				Name:    "flanneld.service",
				Runtime: true,
				DropIns: []config.UnitDropIn{{
					Name: "20-cloudinit.conf",
					Content: `[Service]
Environment="FLANNELD_ETCD_ENDPOINT=http://12.34.56.78:4001"
Environment="FLANNELD_ETCD_PREFIX=/coreos.com/network/tenant1"
`,
				}},
			}}},
		},
	} {
		units := Flannel{tt.config}.Units()
		if !reflect.DeepEqual(units, tt.units) {
			t.Errorf("bad units (%q): want %v, got %v", tt.config, tt.units, units)
		}
	}
}

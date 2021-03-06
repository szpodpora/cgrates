/*
Real-time Online/Offline Charging System (OCS) for Telecom & ISP environments
Copyright (C) ITsysCOM GmbH

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>
*/

package migrator

import (
	"reflect"
	"sort"
	"testing"

	"github.com/cgrates/cgrates/config"
	"github.com/cgrates/cgrates/engine"
	"github.com/cgrates/cgrates/utils"
)

func TestUserProfile2attributeProfile(t *testing.T) {
	usrCfgIn := config.CgrConfig()
	usrCfgIn.MigratorCgrCfg().UsersFilters = []string{"Account"}
	config.SetCgrConfig(usrCfgIn)
	usrTenant := "cgrates.com"
	users := map[int]*v1UserProfile{
		0: &v1UserProfile{
			Tenant:   defaultTenant,
			UserName: "1001",
			Masked:   true,
			Profile:  map[string]string{},
			Weight:   10,
		},
		1: &v1UserProfile{
			Tenant:   usrTenant,
			UserName: "1001",
			Masked:   true,
			Profile: map[string]string{
				"Account": "1002",
				"Subject": "call_1001",
			},
			Weight: 10,
		},
		2: &v1UserProfile{
			Tenant:   defaultTenant,
			UserName: "1001",
			Masked:   false,
			Profile: map[string]string{
				"Account": "1002",
				"ReqType": "*prepaid",
				"msisdn":  "123423534646752",
			},
			Weight: 10,
		},
		3: &v1UserProfile{
			Tenant:   usrTenant,
			UserName: "1001",
			Masked:   false,
			Profile: map[string]string{
				"Account": "1002",
				"ReqType": "*prepaid",
			},
			Weight: 10,
		},
	}
	expected := map[int]*engine.AttributeProfile{
		0: {
			Tenant:             defaultTenant,
			ID:                 "1001",
			Contexts:           []string{utils.META_ANY},
			FilterIDs:          make([]string, 0),
			ActivationInterval: nil,
			Attributes:         []*engine.Attribute{},
			Blocker:            false,
			Weight:             10,
		},
		1: {
			Tenant:             defaultTenant,
			ID:                 "1001",
			Contexts:           []string{utils.META_ANY},
			FilterIDs:          []string{"*string:Account:1002"},
			ActivationInterval: nil,
			Attributes: []*engine.Attribute{
				{
					FieldName:  utils.MetaTenant,
					Initial:    utils.META_ANY,
					Substitute: config.NewRSRParsersMustCompile(usrTenant, true, utils.INFIELD_SEP),
					Append:     true,
				},
				{
					FieldName:  "Subject",
					Initial:    utils.META_ANY,
					Substitute: config.NewRSRParsersMustCompile("call_1001", true, utils.INFIELD_SEP),
					Append:     true,
				},
			},
			Blocker: false,
			Weight:  10,
		},
		2: {
			Tenant:             defaultTenant,
			ID:                 "1001",
			Contexts:           []string{utils.META_ANY},
			FilterIDs:          []string{"*string:Account:1002"},
			ActivationInterval: nil,
			Attributes: []*engine.Attribute{
				{
					FieldName:  utils.RequestType,
					Initial:    utils.META_ANY,
					Substitute: config.NewRSRParsersMustCompile("*prepaid", true, utils.INFIELD_SEP),
					Append:     true,
				},
				{
					FieldName:  "msisdn",
					Initial:    utils.META_ANY,
					Substitute: config.NewRSRParsersMustCompile("123423534646752", true, utils.INFIELD_SEP),
					Append:     true,
				},
			},
			Blocker: false,
			Weight:  10,
		},
		3: {
			Tenant:             defaultTenant,
			ID:                 "1001",
			Contexts:           []string{utils.META_ANY},
			FilterIDs:          []string{"*string:Account:1002"},
			ActivationInterval: nil,
			Attributes: []*engine.Attribute{
				{
					FieldName:  utils.MetaTenant,
					Initial:    utils.META_ANY,
					Substitute: config.NewRSRParsersMustCompile(usrTenant, true, utils.INFIELD_SEP),
					Append:     true,
				},
				{
					FieldName:  utils.RequestType,
					Initial:    utils.META_ANY,
					Substitute: config.NewRSRParsersMustCompile("*prepaid", true, utils.INFIELD_SEP),
					Append:     true,
				},
			},
			Blocker: false,
			Weight:  10,
		},
	}
	for i := range expected {
		rply := userProfile2attributeProfile(users[i])
		sort.Slice(rply.Attributes, func(i, j int) bool {
			if rply.Attributes[i].FieldName == rply.Attributes[j].FieldName {
				return rply.Attributes[i].Initial.(string) < rply.Attributes[j].Initial.(string)
			}
			return rply.Attributes[i].FieldName < rply.Attributes[j].FieldName
		}) // only for test; map returns random keys
		if !reflect.DeepEqual(expected[i], rply) {
			t.Errorf("For %v expected: %s ,received: %s ", i, utils.ToJSON(expected[i]), utils.ToJSON(rply))
		}
	}
}

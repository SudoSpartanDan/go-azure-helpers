package identity

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestLegacySystemUserAssignedListMarshal(t *testing.T) {
	testData := []struct {
		input                           *LegacySystemAndUserAssignedList
		expectedIdentityType            string
		expectedUserAssignedIdentityIds []string
		expectError                     bool
	}{
		{
			input:                           nil,
			expectedIdentityType:            "None",
			expectedUserAssignedIdentityIds: []string{},
		},
		{
			input:                           &LegacySystemAndUserAssignedList{},
			expectedIdentityType:            "None",
			expectedUserAssignedIdentityIds: []string{},
		},
		{
			input: &LegacySystemAndUserAssignedList{
				Type: TypeNone,
			},
			expectedIdentityType:            "None",
			expectedUserAssignedIdentityIds: []string{},
		},
		{
			input: &LegacySystemAndUserAssignedList{
				Type: TypeNone,
				IdentityIds: []string{
					"first",
				},
			},
			expectedIdentityType:            "None",
			expectedUserAssignedIdentityIds: []string{
				// intentionally empty since this is bad data
			},
		},
		{
			input: &LegacySystemAndUserAssignedList{
				Type:        TypeSystemAssigned,
				IdentityIds: []string{},
			},
			expectedIdentityType:            "SystemAssigned",
			expectedUserAssignedIdentityIds: []string{},
		},
		{
			input: &LegacySystemAndUserAssignedList{
				Type:        TypeSystemAssignedUserAssigned,
				IdentityIds: []string{},
			},
			expectedIdentityType:            "SystemAssigned,UserAssigned",
			expectedUserAssignedIdentityIds: []string{},
		},
		{
			input: &LegacySystemAndUserAssignedList{
				Type:        typeLegacySystemAssignedUserAssigned,
				IdentityIds: []string{},
			},
			expectError: true,
		},
		{
			input: &LegacySystemAndUserAssignedList{
				Type:        TypeUserAssigned,
				IdentityIds: []string{},
			},
			expectedIdentityType:            "UserAssigned",
			expectedUserAssignedIdentityIds: []string{},
		},
		{
			input: &LegacySystemAndUserAssignedList{
				Type: TypeSystemAssignedUserAssigned,
				IdentityIds: []string{
					"first",
					"second",
				},
			},
			expectedIdentityType: "SystemAssigned,UserAssigned",
			expectedUserAssignedIdentityIds: []string{
				"first",
				"second",
			},
		},
		{
			input: &LegacySystemAndUserAssignedList{
				Type: typeLegacySystemAssignedUserAssigned,
				IdentityIds: []string{
					"first",
					"second",
				},
			},
			expectError: true,
		},
		{
			input: &LegacySystemAndUserAssignedList{
				Type: TypeUserAssigned,
				IdentityIds: []string{
					"first",
					"second",
				},
			},
			expectedIdentityType: "UserAssigned",
			expectedUserAssignedIdentityIds: []string{
				"first",
				"second",
			},
		},
	}
	for i, v := range testData {
		t.Logf("step %d..", i)

		encoded, err := v.input.MarshalJSON()
		if err != nil {
			if v.expectError {
				continue
			}

			t.Fatalf("encoding: %+v", err)
		}
		if v.expectError {
			t.Fatalf("expected an error but didn't get one")
		}

		var out map[string]interface{}
		if err := json.Unmarshal(encoded, &out); err != nil {
			t.Fatalf("decoding: %+v", err)
		}

		actualIdentityValue := out["type"].(string)
		if v.expectedIdentityType != actualIdentityValue {
			t.Fatalf("expected %q but got %q", v.expectedIdentityType, actualIdentityValue)
		}

		actualUserAssignedIdentityIdsRaw, ok := out["userAssignedIdentities"].([]interface{})
		if !ok {
			if len(v.expectedUserAssignedIdentityIds) == 0 {
				continue
			}

			t.Fatalf("`userAssignedIdentities` was nil")
		}
		actualUserAssignedIdentityIds := make([]string, 0)
		for _, v := range actualUserAssignedIdentityIdsRaw {
			actualUserAssignedIdentityIds = append(actualUserAssignedIdentityIds, v.(string))
		}
		if !reflect.DeepEqual(v.expectedUserAssignedIdentityIds, actualUserAssignedIdentityIds) {
			t.Fatalf("expected %q but got %q", strings.Join(v.expectedUserAssignedIdentityIds, ", "), strings.Join(actualUserAssignedIdentityIds, ", "))
		}
	}
}

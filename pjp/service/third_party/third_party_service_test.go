package thirdparty

import (
	"net/url"
	"scyllax-pjp/model"
	"testing"
)

func TestMasterOutletByIDsEndpointURL_IncludesInactiveFlagWhenEnabled(t *testing.T) {
	url := masterOutletByIDsEndpointURL("http://kong", "1,2,3", 9999, 1)
	expected := "http://kong/v1/outlets?outlet_id=1,2,3&limit=9999&include_inactive=1"
	if url != expected {
		t.Fatalf("expected %s, got %s", expected, url)
	}
}

func TestMasterOutletByIDsEndpointURL_DefaultBehaviorWithoutIncludeInactive(t *testing.T) {
	url := masterOutletByIDsEndpointURL("http://kong", "1,2,3", 9999, 0)
	expected := "http://kong/v1/outlets?outlet_id=1,2,3&limit=9999"
	if url != expected {
		t.Fatalf("expected %s, got %s", expected, url)
	}
}

func TestMasterOutletEndpointURL_ForwardsOutletQueryParams(t *testing.T) {
	endpointURL := masterOutletEndpointURL("http://kong", model.DmsQueryFilter{
		Page:        "2",
		Limit:       "25",
		Query:       "alpha market",
		OutletCode:  "OUT-123",
		OutletID:    123,
		Mode:        "lookup",
		Sort:        "outlet_name:asc",
		IsActive:    "1",
		SalesTeamID: "7",
	})

	parsedURL, err := url.Parse(endpointURL)
	if err != nil {
		t.Fatalf("failed to parse endpoint URL: %v", err)
	}

	expected := map[string]string{
		"limit":               "25",
		"page":                "2",
		"q":                   "alpha market",
		"outlet_code":         "OUT-123",
		"outlet_id":           "123",
		"mode":                "lookup",
		"sort":                "outlet_name:asc",
		"is_active":           "1",
		"sales_team_id":       "7",
		"verification_status": "1",
	}

	query := parsedURL.Query()
	for key, value := range expected {
		if query.Get(key) != value {
			t.Fatalf("expected %s=%s, got %s", key, value, query.Get(key))
		}
	}
}

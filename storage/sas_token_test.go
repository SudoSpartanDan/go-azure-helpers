package storage

import "testing"

func TestParseStorageAccountConnectionString(t *testing.T) {
	testCases := []struct {
		input               string
		expectedAccountName string
		expectedAccountKey  string
		expectedError       bool
	}{
		{
			"DefaultEndpointsProtocol=https;AccountName=azurermtestsa0;AccountKey=2vJrjEyL4re2nxCEg590wJUUC7PiqqrDHjAN5RU304FNUQieiEwS2bfp83O0v28iSfWjvYhkGmjYQAdd9x+6nw==;EndpointSuffix=core.windows.net",
			"azurermtestsa0",
			"2vJrjEyL4re2nxCEg590wJUUC7PiqqrDHjAN5RU304FNUQieiEwS2bfp83O0v28iSfWjvYhkGmjYQAdd9x+6nw==",
			false,
		},
		{
			"DefaultEndpointsProtocol=https;AccountName=azurermtestsa0;AccountKey=2vJrjEyL4re2nxCEg590wJUUC7PiqqrDHjAN5RU304FNUQieiEwS2bfp83O0v28iSfWjvYhkGmjYQAdd9x+6nw==;EndpointSuffix",
			"",
			"",
			true,
		},
	}

	for _, test := range testCases {
		result, err := ParseAccountSASConnectionString(test.input)

		if test.expectedError {
			if err == nil {
				t.Fatalf("Expected error for %s: %q", test.input, err)
			}
			return
		}

		if ! test.expectedError && err != nil {
			t.Fatalf("Failed to parse resource type string: %s, %q", test.input, result)
		}

		if val, pres := result[connStringAccountKeyKey]; !pres || val != test.expectedAccountKey {
			t.Fatalf("Failed to parse Account Key: Expected: %s, Found: %s", test.expectedAccountKey, val)
		}
		if val, pres := result[connStringAccountNameKey]; !pres || val != test.expectedAccountName {
			t.Fatalf("Failed to parse Account Name: Expected: %s, Found: %s", test.expectedAccountName, val)
		}
	}
}

// This connection string was for a real storage account which has been deleted
// so its safe to include here for reference to understand the format.
// DefaultEndpointsProtocol=https;AccountName=azurermtestsa0;AccountKey=T0ZQouXBDpWud/PlTRHIJH2+VUK8D+fnedEynb9Mx638IYnsMUe4mv1fFjC7t0NayTfFAQJzPZuV1WHFKOzGdg==;EndpointSuffix=core.windows.net
func TestComputeSASToken(t *testing.T) {
	testCases := []struct {
		accountName    string
		accountKey     string
		permissions    string
		services       string
		resourceTypes  string
		start          string
		expiry         string
		signedProtocol string
		signedIp       string
		signedVersion  string
		knownSasToken  string
	}{
		{
			"azurermtestsa0",
			"T0ZQouXBDpWud/PlTRHIJH2+VUK8D+fnedEynb9Mx638IYnsMUe4mv1fFjC7t0NayTfFAQJzPZuV1WHFKOzGdg==",
			"rwac",
			"b",
			"c",
			"2018-03-20T04:00:00Z",
			"2020-03-20T04:00:00Z",
			"https",
			"",
			"2017-07-29",
			"?sv=2017-07-29&ss=b&srt=c&sp=rwac&se=2020-03-20T04:00:00Z&st=2018-03-20T04:00:00Z&spr=https&sig=SQigK%2FnFA4pv0F0oMLqr6DxUWV4vtFqWi6q3Mf7o9nY%3D",
		},
		{
			"azurermtestsa0",
			"2vJrjEyL4re2nxCEg590wJUUC7PiqqrDHjAN5RU304FNUQieiEwS2bfp83O0v28iSfWjvYhkGmjYQAdd9x+6nw==",
			"rwdlac",
			"b",
			"sco",
			"2018-03-20T04:00:00Z",
			"2018-03-28T05:04:25Z",
			"https,http",
			"",
			"2017-07-29",
			"?sv=2017-07-29&ss=b&srt=sco&sp=rwdlac&se=2018-03-28T05:04:25Z&st=2018-03-20T04:00:00Z&spr=https,http&sig=OLNwL%2B7gxeDQQaUyNdXcDPK2aCbCMgEkJNjha9te448%3D",
		},
	}

	for _, test := range testCases {
		computedToken, err := ComputeAccountSASToken(test.accountName,
			test.accountKey,
			test.permissions,
			test.services,
			test.resourceTypes,
			test.start,
			test.expiry,
			test.signedProtocol,
			test.signedIp,
			test.signedVersion)

		if err != nil {
			t.Fatalf("Test Failed: Error computing storage account Sas: %q", err)
		}

		if computedToken != test.knownSasToken {
			t.Fatalf("Test failed: Expected Azure SAS %s but was %s", test.knownSasToken, computedToken)
		}
	}
}

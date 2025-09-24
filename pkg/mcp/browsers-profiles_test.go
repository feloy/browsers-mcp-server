package mcp

import (
	"testing"

	"github.com/feloy/browsers-mcp-server/pkg/api"
	"github.com/feloy/browsers-mcp-server/pkg/browsers/test"
	"github.com/google/go-cmp/cmp"
)

func TestBrowsersProfiles(t *testing.T) {
	for _, tt := range []struct {
		name     string
		browsers []api.Browser
		expected []string
	}{
		{
			name: "only one browser with one profile",
			browsers: []api.Browser{
				test.NewBrowser(test.NewBrowserOptions{
					Name:      "browser1",
					Available: true,
					Profiles:  []string{"profile1"},
				}),
			},
			expected: []string{},
		},
		{
			name: "only one browser with multiple profiles",
			browsers: []api.Browser{
				test.NewBrowser(test.NewBrowserOptions{
					Name:      "browser1",
					Available: true,
					Profiles:  []string{"profile1a", "profile1b"},
				}),
			},
			expected: []string{"profile1a", "profile1b"},
		},
		{
			name: "multiple browsers with one profile",
			browsers: []api.Browser{
				test.NewBrowser(test.NewBrowserOptions{
					Name:      "browser1",
					Available: true,
					Profiles:  []string{"profile1"},
				}),
				test.NewBrowser(test.NewBrowserOptions{
					Name:      "browser2",
					Available: true,
					Profiles:  []string{"profile2"},
				}),
			},
			expected: []string{"browser1", "browser2"},
		},
		{
			name: "multiple browsers with one or multiple profiles",
			browsers: []api.Browser{
				test.NewBrowser(test.NewBrowserOptions{
					Name:      "browser1",
					Available: true,
					Profiles:  []string{"profile1a", "profile1b"},
				}),
				test.NewBrowser(test.NewBrowserOptions{
					Name:      "browser2",
					Available: true,
					Profiles:  []string{"profile2"},
				}),
			},
			expected: []string{"profile1a on browser1", "profile1b on browser1", "browser2"},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			browserProfiles := BrowsersProfiles{}
			browserProfiles.Populate(tt.browsers)
			profiles := browserProfiles.FlatList()
			if !cmp.Equal(tt.expected, profiles) {
				t.Errorf("expected %v, got %v", tt.expected, profiles)
			}
		})
	}
}

type testValues struct {
	value           string
	expectedBrowser string
	expectedProfile string
	expectedError   bool
}

func TestGetBrowserAndProfileFromValue(t *testing.T) {
	for _, tt := range []struct {
		name       string
		browsers   []api.Browser
		value      string
		testValues []testValues
	}{
		{
			name: "only one browser with one profile",
			browsers: []api.Browser{
				test.NewBrowser(test.NewBrowserOptions{
					Name:      "browser1",
					Available: true,
					Profiles:  []string{"profile1"},
				}),
			},
			testValues: []testValues{
				{
					value:           "",
					expectedBrowser: "browser1",
					expectedProfile: "profile1",
				},
				// the following ones should not be used by caller, but let support them anyway
				{
					value:           "profile1",
					expectedBrowser: "browser1",
					expectedProfile: "profile1",
				}, {
					value:           "browser1",
					expectedBrowser: "browser1",
					expectedProfile: "profile1",
				}, {
					value:           "profile1 on browser1",
					expectedBrowser: "browser1",
					expectedProfile: "profile1",
				}, {
					value:         "unknown profile on unknown browser",
					expectedError: true,
				},
			},
		},
		{
			name: "only one browser with multiple profiles",
			browsers: []api.Browser{
				test.NewBrowser(test.NewBrowserOptions{
					Name:      "browser1",
					Available: true,
					Profiles:  []string{"profile1a", "profile1b"},
				}),
			},
			testValues: []testValues{
				{
					value:         "",
					expectedError: true,
				},
				{
					value:         "browser1",
					expectedError: true,
				},
				{
					value:           "profile1a",
					expectedBrowser: "browser1",
					expectedProfile: "profile1a",
				},
				{
					value:           "profile1a on browser1",
					expectedBrowser: "browser1",
					expectedProfile: "profile1a",
				},
			},
		},
		{
			name: "multiple browsers with one profile",
			browsers: []api.Browser{
				test.NewBrowser(test.NewBrowserOptions{
					Name:      "browser1",
					Available: true,
					Profiles:  []string{"profile1"},
				}),
				test.NewBrowser(test.NewBrowserOptions{
					Name:      "browser2",
					Available: true,
					Profiles:  []string{"profile2"},
				}),
			},
			testValues: []testValues{
				{
					value:           "browser1",
					expectedBrowser: "browser1",
					expectedProfile: "profile1",
				},
				{
					value:           "browser2",
					expectedBrowser: "browser2",
					expectedProfile: "profile2",
				},
				// Not expected to be used by caller, but let support them anyway
				{
					value:         "",
					expectedError: true,
				},
				{
					value:         "profile1",
					expectedError: true,
				},
				{
					value:         "profile2",
					expectedError: true,
				},
				{
					value:           "profile1 on browser1",
					expectedBrowser: "browser1",
					expectedProfile: "profile1",
				},
				{
					value:           "profile2 on browser2",
					expectedBrowser: "browser2",
					expectedProfile: "profile2",
				},
			},
		},
		{
			name: "multiple browsers with one or multiple profiles",
			browsers: []api.Browser{
				test.NewBrowser(test.NewBrowserOptions{
					Name:      "browser1",
					Available: true,
					Profiles:  []string{"profile1a", "profile1b"},
				}),
				test.NewBrowser(test.NewBrowserOptions{
					Name:      "browser2",
					Available: true,
					Profiles:  []string{"profile2"},
				}),
			},
			testValues: []testValues{
				{
					value:           "browser2",
					expectedBrowser: "browser2",
					expectedProfile: "profile2",
				},
				{
					value:           "profile1a on browser1",
					expectedBrowser: "browser1",
					expectedProfile: "profile1a",
				},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			for _, testValue := range tt.testValues {
				browser, profile, err := GetBrowserAndProfileFromValue(testValue.value, tt.browsers)
				if testValue.expectedError && err == nil {
					t.Errorf("value %q: expected error, got nil", testValue.value)
				}
				if !testValue.expectedError && err != nil {
					t.Errorf("value %q: expected no error, got %v", testValue.value, err)
				}
				if testValue.expectedBrowser != browser {
					t.Errorf("value %q: expected browser %q, got %q", testValue.value, testValue.expectedBrowser, browser)
				}
				if testValue.expectedProfile != profile {
					t.Errorf("value %q: expected profile %q, got %q", testValue.value, testValue.expectedProfile, profile)
				}
			}
		})
	}
}

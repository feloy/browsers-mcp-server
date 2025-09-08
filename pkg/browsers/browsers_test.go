package browsers

import (
	"errors"
	"testing"

	"github.com/feloy/mcp-server/pkg/browsers/test"
)

func TestGetBrowsers(t *testing.T) {
	for _, tt := range []struct {
		name     string
		browsers []*test.Browser
		expected []string
	}{
		{
			name: "only available browsers are returned",
			browsers: []*test.Browser{
				test.NewBrowser("browser1", true, nil, []string{"profile1"}, nil),
				test.NewBrowser("browser2", false, nil, []string{"profile2"}, nil),
			},
			expected: []string{"browser1"},
		},
		{
			name: "available error is not returned",
			browsers: []*test.Browser{
				test.NewBrowser("browser1", false, errors.New("an errror"), []string{"profile1"}, nil),
				test.NewBrowser("browser2", true, nil, []string{"profile2"}, nil),
			},
			expected: []string{"browser2"},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			Clear()
			for _, browser := range tt.browsers {
				Register(browser)
			}
			browsers := GetBrowsers()
			if len(browsers) != len(tt.expected) {
				t.Errorf("Expected %d browsers, got %d", len(tt.expected), len(browsers))
			}
			for i, browser := range browsers {
				if browser.Name() != tt.expected[i] {
					t.Errorf("Expected browser %s, got %s", tt.expected[i], browser.Name())
				}
			}
		})
	}
}

func TestGetBrowserByName(t *testing.T) {
	for _, tt := range []struct {
		name          string
		browserName   string
		browsers      []*test.Browser
		expected      string
		expectedError error
	}{
		{
			name:        "browser is found",
			browserName: "browser1",
			browsers: []*test.Browser{
				test.NewBrowser("browser1", true, nil, []string{"profile1"}, nil),
			},
			expected:      "browser1",
			expectedError: nil,
		},
		{
			name:        "browser not found",
			browserName: "browser2",
			browsers: []*test.Browser{
				test.NewBrowser("browser1", true, nil, []string{"profile1"}, nil),
			},
			expected:      "",
			expectedError: errors.New("browser \"browser2\" not found"),
		},
		{
			name:        "browser is not available",
			browserName: "browser1",
			browsers: []*test.Browser{
				test.NewBrowser("browser1", false, nil, []string{"profile1"}, nil),
			},
			expected:      "",
			expectedError: errors.New("browser \"browser1\" is not available"),
		},
		{
			name:        "browser is available error",
			browserName: "browser1",
			browsers: []*test.Browser{
				test.NewBrowser("browser1", true, errors.New("an error"), []string{"profile1"}, nil),
			},
			expected:      "",
			expectedError: errors.New("an error"),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			Clear()
			for _, browser := range tt.browsers {
				Register(browser)
			}
			browser, err := GetBrowserByName(tt.browserName)
			if err != nil && err.Error() != tt.expectedError.Error() {
				t.Errorf("Expected error %v, got %v", tt.expectedError, err)
			}
			if err == nil && tt.expectedError != nil {
				t.Errorf("Expected error %v, got nil", tt.expectedError)
			}
			if tt.expectedError == nil {
				if browser.Name() != tt.expected {
					t.Errorf("Expected browser %q, got %q", tt.expected, browser.Name())
				}
			}
		})
	}
}

package created_by

import (
	"net/http"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/tests"
)

const testDataDir = "./test/pb_data/"

func generateToken(collectionNameOrId string, email string) (string, string, error) {
	app, err := tests.NewTestApp(testDataDir)
	if err != nil {
		return "", "", err
	}
	defer app.Cleanup()

	record, err := app.FindAuthRecordByEmail(collectionNameOrId, email)
	if err != nil {
		return "", "", err
	}
	token, err := record.NewAuthToken()
	if err != nil {
		return "", "", err
	}

	return token, record.Id, nil
}

func TestPlugin(t *testing.T) {
	setupTestApp := func(t testing.TB) *tests.TestApp {
		testApp, err := tests.NewTestApp(testDataDir)
		if err != nil {
			t.Fatal(err)
		}
		(&Plugin{
			// test config will go here
			Fields: []string{
				"posts.user_id",
			},
		}).Init(testApp)
		return testApp
	}

	userToken, userId, err := generateToken("users", "test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []tests.ApiScenario{
		{
			Name:   "create record",
			Method: http.MethodPost,
			URL:    "/api/collections/posts/records",
			Headers: map[string]string{
				"Authorization": userToken,
			},
			Body:           strings.NewReader(`{}`),
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"user_id":"` + userId + `"`,
			},
			TestAppFactory: setupTestApp,
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestValidation(t *testing.T) {
	type scenario struct {
		name      string
		input     *Plugin
		expectErr bool
	}

	scenarios := []scenario{
		{
			"empty",
			&Plugin{},
			false,
		},
		{
			"valid",
			&Plugin{
				Fields: []string{
					"collection_name.field_name",
				},
			},
			false,
		},
		{
			"no_period",
			&Plugin{
				Fields: []string{
					"no_period",
				},
			},
			true,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			err := s.input.Validate()
			if s.expectErr && err == nil {
				t.Error("expected error and got none")
			}
			if !s.expectErr && err != nil {
				t.Errorf("unexpected error: %s", err)
			}
		})
	}
}

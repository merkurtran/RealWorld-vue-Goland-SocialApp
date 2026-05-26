package tests

import (
	"Server/models"
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRegistration(t *testing.T) {
	// clear up collection before the test
	cleanupCollection()

	tests := []struct {
		name           string
		payload        models.CreateUser
		expectedStatus int
		shouldContain  []string
	}{
		{
			name: "Valid Registration",
			payload: models.CreateUser{
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
			},
			expectedStatus: 201,
			shouldContain:  []string{"result", "token"},
		},
		{
			name: "Missing Required Fields",
			payload: models.CreateUser{
				Email: "missing@example.com",
			},
			expectedStatus: 400,
		},
		{
			name: "Duplicate Email Registration",
			payload: models.CreateUser{
				Email:     "duplicate@example.com",
				Password:  "password456",
				FirstName: "jane",
				LastName:  "smith",
			},
			expectedStatus: 400,
			shouldContain:  []string{"Already Exit!"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// for duplicate test we meed first to register the original user
			if tt.name == "Duplicate Email Registration" {
				firstpayload := models.CreateUser{
					Email:     "duplicate@example.com",
					Password:  "password2226",
					FirstName: "ahmed",
					LastName:  "khalaf",
				}
				registerUser(t, firstpayload, 201)
			}

			// making the request
			jsonPayload, err := json.Marshal(tt.payload)
			require.NoError(t, err)

			req, err := http.NewRequest("POST", "/user/signup", bytes.NewBuffer(jsonPayload))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			// check Status code
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			// parse res
			var responseBody any
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			if err != nil {
				t.Logf("could not parse response as json: %v", err)
				return
			}

			// convert to string for contains check
			responseStr, _ := json.Marshal(responseBody)

			// check expected content
			for _, contain := range tt.shouldContain {
				assert.Contains(t, string(responseStr), contain)
			}

			// For successful registration verfify basic structer
			if tt.expectedStatus == 201 {
				if respMap, ok := responseBody.(map[string]any); ok {
					assert.Contains(t, respMap, "token")
					assert.Contains(t, respMap, "result")

					if result, ok := respMap["result"].(map[string]any); ok {
						assert.Equal(t, tt.payload.Email, result["email"])
						expactedName := tt.payload.FirstName + " " + tt.payload.LastName
						assert.Equal(t, expactedName, result["name"])
					}
				}
			}

		})
	}
}

// helper func to create user for testing
func registerUser(t *testing.T, payload models.CreateUser, expectedStatus int) {
	jsonPayload, err := json.Marshal(payload)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/user/signup", bytes.NewBuffer(jsonPayload))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, expectedStatus, resp.StatusCode)
}

func TestUserLogin(t *testing.T) {
	cleanupCollection()

	// register a user first
	registerPayload := models.CreateUser{
		Email:     "login@example.com",
		Password:  "password123",
		FirstName: "Login",
		LastName:  "Test",
	}
	registerUser(t, registerPayload, 201)

	tests := []struct {
		name           string
		payload        models.LoginUser
		expectedStatus int
		shouldContain  []string
	}{
		{
			name: "Valid Login",
			payload: models.LoginUser{
				Email:    "login@example.com",
				Password: "password123",
			},
			expectedStatus: 200,
			shouldContain:  []string{"token", "result"},
		},
		{
			name: "Invalid Email",
			payload: models.LoginUser{
				Email:    "noneexistent@example.com",
				Password: "password123",
			},
			expectedStatus: 400,
		},
		{
			name: "Invalid Password",
			payload: models.LoginUser{
				Email:    "login@example.com",
				Password: "wrongpassword",
			},
			expectedStatus: 400,
		},
		{
			name: "Empty Credentials",
			payload: models.LoginUser{
				Email:    "",
				Password: "",
			},
			expectedStatus: 400,
		},
	}
	// loop for case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// making the request
			jsonPayload, err := json.Marshal(tt.payload)
			require.NoError(t, err)

			req, err := http.NewRequest("POST", "/user/signin", bytes.NewBuffer(jsonPayload))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			// check status code
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			// parse res
			var responseBody any
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			if err != nil {
				t.Logf("could not parse response as json: %v", err)
				return
			}

			// convert to string for contains check
			responseStr, _ := json.Marshal(responseBody)

			// check expected content
			for _, contain := range tt.shouldContain {
				assert.Contains(t, string(responseStr), contain)
			}

			// for successful login verify base struct
			if tt.expectedStatus == 200 {
				if respMap, ok := responseBody.(map[string]any); ok {
					assert.Contains(t, respMap, "token")
					assert.Contains(t, respMap, "result")

					if result, ok := respMap["result"].(map[string]any); ok {
						assert.Equal(t, tt.payload.Email, result["email"])
					}
				}
			}
		})
	}

}

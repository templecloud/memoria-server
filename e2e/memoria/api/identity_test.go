package identity

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/templecloud/memoria-server/internal/memoria/boot"
	"github.com/templecloud/memoria-server/internal/memoria/controller/identity"

	"github.com/templecloud/memoria-server/e2e/memoria/framework"
)

// TestIdentity tests the basic sequence:
//
// 1. Not Logged in. '/api/v1/health' endpoint returns 401.
// 2/ Login successfully
// 3. Not Logged in. '/api/v1/health' endpoint returns 200.
//
// NB: Still currently need to initialise MongoDB and ensure the User exists before executing the test.
//
func TestIdentity(t *testing.T) {

	fw := framework.NewFramework("")
	assert.NotNil(t, fw)
	fw.BeforeEach()

	router := boot.NewServer()

	// curl -v --cookie "token=${JWT}" localhost:8080/api/v1/health
	actual := invoke(router, "GET", "/api/v1/health", nil, nil)
	assert.NotNil(t, actual)
	assert.Equal(t, http.StatusUnauthorized, actual.Code)

	// curl -v -X POST localhost:8080/api/v1/signup -d '{ "name": "test-user", "email": "test@test.com", "password": "test" }
	login := identity.Login{Email: "test@test.com", Password: "test"}
	signup := identity.Signup{Name: "Test", Login: login}
	body, _ := json.Marshal(signup)
	reader := bytes.NewReader(body)
	actual = invoke(router, "POST", "/api/v1/signup", reader, nil)
	
	assert.Equal(t, http.StatusOK, actual.Code)

   	// curl -v localhost:8080/api/v1/login -d '{ "email": "test@test.com", "password": "test" }'	
	body, _ = json.Marshal(login)
	reader = bytes.NewReader(body)
	actual = invoke(router, "POST", "/api/v1/login", reader, nil)
	assert.NotNil(t, actual)
	assert.Equal(t, http.StatusOK, actual.Code)

   	// curl -v --cookie "token=${JWT}" localhost:8080/api/v1/health
	cookies := actual.Result().Cookies()
	actual = invoke(router, "GET", "/api/v1/health", nil, cookies)
	assert.NotNil(t, actual)
	assert.Equal(t, http.StatusOK, actual.Code)

	fw.AfterEach()
}

func invoke(
	handler http.Handler, method,
	path string,
	body io.Reader,
	cookies []*http.Cookie,
) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	if cookies != nil {
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
	}
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)
	return recorder
}
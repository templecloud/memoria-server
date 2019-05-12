package identity_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/templecloud/memoria-server/e2e/memoria/framework"
	"github.com/templecloud/memoria-server/internal/memoria/boot"
	"github.com/templecloud/memoria-server/internal/memoria/controller/identity"
)

func TestIdentity(t *testing.T) {
	RegisterFailHandler(Fail)          // Register `gomega` with `ginkgo`.
	RunSpecs(t, "Identity Test Suite") // Run `ginkgo`.
}

var fw *framework.Framework

var _ = BeforeSuite(func() {
	fw = framework.NewFramework()
	fw.BeforeEach()
})

var _ = AfterSuite(func() {
	fw.AfterEach()
})

var _ = Describe("Identity", func() {

	Describe("Checking user login", func() {

		Context("When not logged in", func() {

			It("should not be possible to get the health of the server", func() {
				router := boot.NewServer()
				// curl -v localhost:8080/api/v1/health
				actual := invoke(router, "GET", "/api/v1/health", nil, nil)
				Expect(actual).NotTo(BeNil())
				Expect(actual.Code).To(Equal(http.StatusUnauthorized))
			})
		})

		Context("When logged in", func() {

			It("should not be possible to get the health of the server", func() {
				router := boot.NewServer()
				login := identity.Login{Email: "test@test.com", Password: "test"}
				signup := identity.Signup{Name: "Test", Login: login}

				// curl -v -X POST localhost:8080/api/v1/signup -d '{ "name": "test-user", "email": "test@test.com", "password": "test" }
				body, _ := json.Marshal(signup)
				reader := bytes.NewReader(body)
				actual := invoke(router, "POST", "/api/v1/signup", reader, nil)
				Expect(actual).NotTo(BeNil())
				Expect(actual.Code).To(Equal(http.StatusOK))

				// curl -v localhost:8080/api/v1/login -d '{ "email": "test@test.com", "password": "test" }'
				body, _ = json.Marshal(login)
				reader = bytes.NewReader(body)
				actual = invoke(router, "POST", "/api/v1/login", reader, nil)
				Expect(actual).NotTo(BeNil())
				Expect(actual.Code).To(Equal(http.StatusOK))
				cookies := actual.Result().Cookies()
				Expect(cookies).NotTo(BeNil())

				// curl -v --cookie "token=${JWT}" localhost:8080/api/v1/health
				actual = invoke(router, "GET", "/api/v1/health", nil, cookies)
				Expect(actual).NotTo(BeNil())
				Expect(actual.Code).To(Equal(http.StatusOK))
			})
		})
	})
})

// invoke the specified server endpoint and record the results.
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

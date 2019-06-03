package identity_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/templecloud/memoria-server/internal/memoria/boot"
	"github.com/templecloud/memoria-server/internal/memoria/controller/identity"
	"github.com/templecloud/memoria-server/test/e2e/memoria/framework"
)

// curl -v localhost:8080/api/v1/health
// curl -v -X POST localhost:8080/api/v1/signup -d '{ "name": "test-user", "email": "test@test.com", "password": "test" }
// curl -v localhost:8080/api/v1/login -d '{ "email": "test@test.com", "password": "test" }'
// curl -v --cookie "token=${JWT}" localhost:8080/api/v1/health

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
				router := boot.NewDefaultServer()
				actual := invoke(router, "GET", "/api/v1/health", nil, nil)
				Expect(actual).NotTo(BeNil())
				Expect(actual.Code).To(Equal(http.StatusUnauthorized))
			})
		})

		Context("When logged in", func() {

			It("should not be possible to get the health of the server", func() {
				router := boot.NewDefaultServer()
				login := identity.Login{Email: "test@test.com", Password: "test"}
				signup := identity.Signup{Name: "Test", Login: login}

				actual := invoke(router, "POST", "/api/v1/signup", signup, nil)
				Expect(actual).NotTo(BeNil())
				Expect(actual.Code).To(Equal(http.StatusOK))

				actual = invoke(router, "POST", "/api/v1/login", login, nil)
				Expect(actual).NotTo(BeNil())
				Expect(actual.Code).To(Equal(http.StatusOK))
				cookies := actual.Result().Cookies()
				Expect(cookies).NotTo(BeNil())

				actual = invoke(router, "GET", "/api/v1/health", nil, cookies)
				Expect(actual).NotTo(BeNil())
				Expect(actual.Code).To(Equal(http.StatusOK))
			})
		})
	})

	Describe("Checking user creation", func() {

		Context("When not logged in", func() {

			login := identity.Login{Email: "test2@test.com", Password: "test2"}
			signup := identity.Signup{Name: "Test2", Login: login}

			It("should be possible to create a new user", func() {
				router := boot.NewDefaultServer()

				actual := invoke(router, "POST", "/api/v1/signup", signup, nil)
				Expect(actual).NotTo(BeNil())
				Expect(actual.Code).To(Equal(http.StatusOK))
			})

			It("should not be possible to create a new user when the username is taken", func() {
				router := boot.NewDefaultServer()

				actual := invoke(router, "POST", "/api/v1/signup", signup, nil)
				Expect(actual).NotTo(BeNil())
				Expect(actual.Code).To(Equal(http.StatusConflict))
				Expect(body(actual)["errorMessage"]).To(Equal("User already registered."))
			})
		})
	})

})

// invoke the specified server endpoint and record the results.
func invoke(
	handler http.Handler, method,
	path string,
	data interface{},
	cookies []*http.Cookie,
) *httptest.ResponseRecorder {
	body, _ := json.Marshal(data)
	reader := bytes.NewReader(body)
	req, _ := http.NewRequest(method, path, reader)
	if cookies != nil {
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
	}
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)
	return recorder
}

func body(recorded *httptest.ResponseRecorder) map[string]interface{} {
	var body map[string]interface{}
	err := json.Unmarshal(recorded.Body.Bytes(), &body)
	if err != nil {
		panic(err)
	}
	return body
}


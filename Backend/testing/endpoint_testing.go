package testing

import (
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
	_ "github.com/onsi/ginkgo/v2"
	_ "github.com/onsi/gomega"
)

var _ = Describe("User API", func() {
	var (
		e   *echo.Echo
		req *http.Request
		rec *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		e = echo.New()
		e.GET("/users/name/:name", func(c echo.Context) error {
			name := c.Param("name")
			if name == "alice" {
				return c.JSON(http.StatusOK, map[string]string{"name": "alice"})
			}
			return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
		})
	})

	It("returns user when name exists", func() {
		req = httptest.NewRequest(http.MethodGet, "/users/name/alice", nil)
		rec = httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		Expect(rec.Code).To(Equal(http.StatusOK))
		Expect(rec.Body.String()).To(ContainSubstring("alice"))
	})

	It("returns 404 when user does not exist", func() {
		req = httptest.NewRequest(http.MethodGet, "/users/name/bob", nil)
		rec = httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		Expect(rec.Code).To(Equal(http.StatusNotFound))
	})
})

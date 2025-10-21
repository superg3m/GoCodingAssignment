package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/superg3m/server/Model"
)

func TestAPI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "API Integration Suite")
}

var baseURL string

var _ = BeforeSuite(func() {
	baseURL = "http://localhost:8080"
})

func sendRequest(method, url string, body any) (*http.Response, error) {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, url, &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	return client.Do(req)
}

func decodeResponseBody(resp *http.Response, out any) {
	data, err := io.ReadAll(resp.Body)
	Expect(err).To(BeNil())
	if out != nil {
		Expect(json.Unmarshal(data, out)).To(Succeed())
	}
}

var _ = Describe("User CRUD API)", func() {
	It("creates, updates, and deletes a user safely", func() {
		createBody := map[string]any{
			"user_name":   "bob_test",
			"first_name":  "Bob",
			"last_name":   "Smith",
			"email":       "bob_test@example.com",
			"user_status": "active",
			"department":  "IT",
		}

		var u Model.User
		cres, err := sendRequest(http.MethodPost, baseURL+"/User/Create", createBody)
		Expect(err).To(BeNil())
		defer cres.Body.Close()
		Expect(cres.StatusCode).To(Equal(http.StatusOK))
		decodeResponseBody(cres, &u)

		DeferCleanup(func() {
			deleteBody := map[string]any{"user_id": u.ID}
			dres, err := sendRequest(http.MethodDelete, baseURL+"/User/Delete", deleteBody)
			Expect(err).To(BeNil())
			if dres != nil {
				defer dres.Body.Close()
				Expect(dres.StatusCode).To(Equal(http.StatusOK))
			}
		})

		updateBody := map[string]any{
			"user_id":     u.ID,
			"user_name":   "bob_test_new",
			"first_name":  "Bob",
			"last_name":   "Smith",
			"email":       "bob_new@example.com",
			"user_status": "active",
			"department":  "IT",
		}
		ures, err := sendRequest(http.MethodPatch, baseURL+"/User/Update", updateBody)
		Expect(err).To(BeNil())
		defer ures.Body.Close()
		Expect(ures.StatusCode).To(Equal(http.StatusOK))
	})

	It("rejects invalid email formats", func() {
		body := map[string]any{
			"user_name":   "invalid_email_user",
			"first_name":  "Test",
			"last_name":   "User",
			"email":       "not-an-email",
			"user_status": "active",
			"department":  "IT",
		}
		resp, err := sendRequest(http.MethodPost, baseURL+"/User/Create", body)
		Expect(err).To(BeNil())
		defer resp.Body.Close()
		Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
	})

	It("rejects invalid id for delete request", func() {
		body := map[string]any{
			"user_id": 0,
		}
		resp, err := sendRequest(http.MethodDelete, baseURL+"/User/Delete", body)
		Expect(err).To(BeNil())
		defer resp.Body.Close()
		Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
	})

	It("rejects duplicate usernames", func() {
		userBody := map[string]any{
			"user_name":   "duplicate_user",
			"first_name":  "John",
			"last_name":   "Doe",
			"email":       "dupe1@example.com",
			"user_status": "active",
			"department":  "Sales",
		}

		var u1 Model.User
		r1, err := sendRequest(http.MethodPost, baseURL+"/User/Create", userBody)
		Expect(err).To(BeNil())
		defer r1.Body.Close()
		Expect(r1.StatusCode).To(Equal(http.StatusOK))
		decodeResponseBody(r1, &u1)

		DeferCleanup(func() {
			deleteBody := map[string]any{"user_id": u1.ID}
			dres, _ := sendRequest(http.MethodDelete, baseURL+"/User/Delete", deleteBody)
			if dres != nil {
				defer dres.Body.Close()
			}
		})

		userBody["email"] = "dupe2@example.com"
		r2, err := sendRequest(http.MethodPost, baseURL+"/User/Create", userBody)
		Expect(err).To(BeNil())
		defer r2.Body.Close()
		Expect(r2.StatusCode).To(Equal(http.StatusConflict))
	})

	It("rejects missing required fields", func() {
		body := map[string]any{
			"user_name": "missing_email",
		}
		resp, err := sendRequest(http.MethodPost, baseURL+"/User/Create", body)
		Expect(err).To(BeNil())
		defer resp.Body.Close()
		Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
	})
})

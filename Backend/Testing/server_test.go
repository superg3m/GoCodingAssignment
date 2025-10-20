package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/superg3m/server/Model"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAPI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "API Integration Suite")
}

var baseURL string

var _ = BeforeSuite(func() {
	baseURL = "http://localhost:8080" // your running server
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
	bodyBytes, err := io.ReadAll(resp.Body)
	Expect(err).To(BeNil())

	err = json.Unmarshal(bodyBytes, out)
	Expect(err).To(BeNil())
}

var _ = Describe("User API (live server)", func() {
	It("creates, updates, and deletes a user successfully", func() {
		createBody := map[string]any{
			"user_name":   "bob",
			"first_name":  "Bob",
			"last_name":   "Smith",
			"email":       "bob@example.com",
			"user_status": "active",
			"department":  "IT",
		}
		var u Model.User
		cres, err := sendRequest(http.MethodPost, baseURL+"/User/Create", createBody)
		Expect(err).To(BeNil())
		defer cres.Body.Close()
		Expect(cres.StatusCode).To(Equal(http.StatusOK))
		decodeResponseBody(cres, &u)

		updateBody := map[string]any{
			"user_id":     u.ID,
			"user_name":   "Newbob",
			"first_name":  "Bob",
			"last_name":   "Smith",
			"email":       "testing@example.com",
			"user_status": "active",
			"department":  "IT",
		}
		ures, err := sendRequest(http.MethodPatch, baseURL+"/User/Update", updateBody)
		Expect(err).To(BeNil())
		defer ures.Body.Close()
		Expect(ures.StatusCode).To(Equal(http.StatusOK))

		deleteBody := map[string]any{"user_id": u.ID}
		dres, err := sendRequest(http.MethodDelete, baseURL+"/User/Delete", deleteBody)
		Expect(err).To(BeNil())
		defer dres.Body.Close()
		Expect(dres.StatusCode).To(Equal(http.StatusOK))
	})
})

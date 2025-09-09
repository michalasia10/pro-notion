package server

import (
	"io"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("HelloWorldHandler", func() {
	It("returns 200 and JSON body", func() {
		s := &Server{}
		ts := httptest.NewServer(http.HandlerFunc(s.HelloWorldHandler))
		DeferCleanup(ts.Close)

		resp, err := http.Get(ts.URL)
		Expect(err).ToNot(HaveOccurred())
		DeferCleanup(func() { _ = resp.Body.Close() })

		Expect(resp.StatusCode).To(Equal(http.StatusOK))
		expected := "{\"message\":\"Hello World\"}"
		body, err := io.ReadAll(resp.Body)
		Expect(err).ToNot(HaveOccurred())
		Expect(string(body)).To(Equal(expected))
	})
})

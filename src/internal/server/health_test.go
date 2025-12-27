package server

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type mockDB struct {
	status string
}

func (m mockDB) Health() map[string]string {
	return map[string]string{
		"status":  m.status,
		"message": "",
	}
}

func (mockDB) Close() error { return nil }

var _ = Describe("healthHandler", func() {
	It("returns 503 when DB is down", func() {
		s := &Server{
			db:          mockDB{status: "down"},
			redisClient: nil,
		}

		req := httptest.NewRequest("GET", "/health", nil)
		rec := httptest.NewRecorder()

		s.healthHandler(rec, req)

		Expect(rec.Code).To(Equal(http.StatusServiceUnavailable))
		Expect(rec.Body.String()).To(ContainSubstring("down"))
	})

	It("returns 503 when Redis is down", func() {
		s := &Server{
			db:          mockDB{status: "up"},
			redisClient: nil,
		}

		req := httptest.NewRequest("GET", "/health", nil)
		rec := httptest.NewRecorder()

		s.healthHandler(rec, req)

		Expect(rec.Code).To(Equal(http.StatusServiceUnavailable))
		Expect(rec.Body.String()).To(ContainSubstring("Redis"))
	})
})

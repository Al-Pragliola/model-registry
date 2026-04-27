package api

import (
	"encoding/csv"
	"net/http"
	"strings"

	"github.com/kubeflow/hub/ui/bff/internal/integrations/kubernetes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ExportCatalogModelsHandler", func() {
	Context("testing Export Catalog Models Handler", Ordered, func() {

		It("should export all models as CSV", func() {
			By("fetching CSV export for all models")
			requestIdentity := kubernetes.RequestIdentity{
				UserID: "user@example.com",
			}

			// Use setupApiTest with string type to get raw body
			// setupApiTest will fail JSON unmarshal, so we use a direct approach
			raw, rs, err := setupRawApiTest(
				http.MethodGet,
				"/api/v1/model_catalog/models/export?namespace=kubeflow",
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"kubeflow",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusOK))

			By("checking response headers")
			Expect(rs.Header.Get("Content-Type")).To(Equal("text/csv; charset=utf-8"))
			Expect(rs.Header.Get("Content-Disposition")).To(ContainSubstring("attachment; filename="))
			Expect(rs.Header.Get("Content-Disposition")).To(ContainSubstring("model-catalog-export-"))

			By("parsing CSV content")
			// Strip BOM if present
			body := strings.TrimPrefix(raw, "\xEF\xBB\xBF")
			reader := csv.NewReader(strings.NewReader(body))
			records, err := reader.ReadAll()
			Expect(err).NotTo(HaveOccurred())

			By("checking CSV has header row and data rows")
			Expect(len(records)).To(BeNumerically(">", 1))

			headers := records[0]
			Expect(headers).To(ContainElements("Name", "Source ID", "Provider", "Description"))
			Expect(headers).To(ContainElements("Language", "Tasks", "Created", "Last Updated"))

			By("checking data rows have correct number of columns")
			for _, row := range records[1:] {
				Expect(len(row)).To(Equal(len(headers)))
			}

			By("checking model names are present")
			names := make([]string, 0)
			nameIdx := 0
			for i, h := range headers {
				if h == "Name" {
					nameIdx = i
					break
				}
			}
			for _, row := range records[1:] {
				names = append(names, row[nameIdx])
			}
			Expect(names).NotTo(BeEmpty())
		})

		It("should export filtered models as CSV", func() {
			By("fetching CSV export with source filter")
			requestIdentity := kubernetes.RequestIdentity{
				UserID: "user@example.com",
			}

			raw, rs, err := setupRawApiTest(
				http.MethodGet,
				"/api/v1/model_catalog/models/export?namespace=kubeflow&source=sample-source",
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"kubeflow",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusOK))

			body := strings.TrimPrefix(raw, "\xEF\xBB\xBF")
			reader := csv.NewReader(strings.NewReader(body))
			records, err := reader.ReadAll()
			Expect(err).NotTo(HaveOccurred())

			By("checking all returned models belong to filtered source")
			headers := records[0]
			sourceIdx := -1
			for i, h := range headers {
				if h == "Source ID" {
					sourceIdx = i
					break
				}
			}
			Expect(sourceIdx).NotTo(Equal(-1))

			for _, row := range records[1:] {
				Expect(row[sourceIdx]).To(Equal("sample-source"))
			}
		})

		It("should return headers-only CSV when no models match", func() {
			By("fetching CSV export with non-matching source")
			requestIdentity := kubernetes.RequestIdentity{
				UserID: "user@example.com",
			}

			raw, rs, err := setupRawApiTest(
				http.MethodGet,
				"/api/v1/model_catalog/models/export?namespace=kubeflow&source=nonexistent-source",
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"kubeflow",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusOK))

			body := strings.TrimPrefix(raw, "\xEF\xBB\xBF")
			reader := csv.NewReader(strings.NewReader(body))
			records, err := reader.ReadAll()
			Expect(err).NotTo(HaveOccurred())

			By("checking only header row is present")
			Expect(len(records)).To(Equal(1))
			Expect(records[0]).To(ContainElements("Name", "Source ID", "Provider"))
		})
	})
})

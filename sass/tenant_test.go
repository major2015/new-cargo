package sass

import (
	_ "fmt" // for adhoc printing

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	qm "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var _ = Describe("Tenant methods", func() {
	Test("can query", &TestFlags{}, func(env *TestEnv) {
		Expect(
			Tenants(qm.Where("id=?", env.Tenant.ID)).ExistsP(env.DB),
		).To(Equal(true))
	})
})

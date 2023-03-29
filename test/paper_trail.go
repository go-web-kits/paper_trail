package test

import (
	"github.com/go-web-kits/dbx"
	. "github.com/go-web-kits/paper_trail"
	. "github.com/go-web-kits/testx"
	"github.com/go-web-kits/testx/factory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PaperTrail", func() {
	var (
		paper   Paper
		version Version
		p       *MonkeyPatches
	)

	var FindVersion = func() bool {
		return Expect(dbx.FindBy(&version, dbx.EQ{"item_id": paper.ID}, dbx.With{Order: "id DESC"})).To(HaveFound())

	}

	BeforeEach(func() {
		paper = Paper{Title: "test"}
	})

	AfterEach(func() {
		CleanData(&Paper{}, &Version{})
		Reset(&paper, &version)
		p.Check()
	})

	Describe("Trailable.Trail", func() {
		It("creates a version by given changes", func() {
			factory.Create(&paper)
			Expect(Paper{}.Trail(dbx.Conn().DB, &paper, h{"xyz": []int{1, 123}}, "will")).To(Succeed())
			FindVersion()
			Expect(version.Event).To(Equal("update xyz"))
			Expect(version.Whodunnit).To(Equal("will"))
			Expect(version.Object).To(EqualYAML(h{"jsonb": "", "nilable": nil, "title": "test"}))
			Expect(version.ObjectChanges).To(EqualYAML(h{"xyz": []int{1, 123}}))
		})
	})

	Describe("GetVersionsOf", func() {
		It("returns a list of versions which are of the paper", func() {
			factory.Create(&paper)
			Expect(Paper{}.Trail(dbx.Conn().DB, &paper, h{"xyz": []int{1, 123}}, "will")).To(Succeed())
			Expect(GetVersionsOf(&paper)).To(HaveLen(2))
			// []h{
			// 				{
			// 					"event":   "create",
			// 					"admin":   "",
			// 					"changes": h{"title": []interface{}{nil, "test"}},
			// 					"time":    "2019-09-24T20:04:29.187464+08:00",
			// 				},
			// 				{
			// 					"changes": h{"xyz": []interface{}{1, 123}},
			// 					"time":    "2019-09-24T20:04:29.18839+08:00",
			// 					"event":   "update xyz",
			// 					"admin":   "will",
			// 				},
			// 			}
		})
	})
})

package test

import (
	"github.com/gin-gonic/gin"
	"github.com/go-web-kits/dbx"
	. "github.com/go-web-kits/paper_trail"
	. "github.com/go-web-kits/testx"
	"github.com/go-web-kits/testx/factory"
	"github.com/go-web-kits/utils/structx"
	"github.com/k0kubun/pp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Version", func() {
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

	Describe("Create Version via GORM callbacks", func() {
		When("the paper is being created", func() {
			It("creates the version which event is create", func() {
				factory.Create(&paper)
				FindVersion()
				Expect(version).To(HaveAttributes(Version{Event: "create", ItemID: paper.ID, ItemType: "Paper"}, "time"))
				Expect(version.Object).To(EqualYAML(h{"jsonb": "", "nilable": nil, "title": "test"}))
				Expect(version.ObjectChanges).To(EqualYAML(h{"jsonb": []interface{}{nil, ""}, "title": []interface{}{nil, "test"}}))
			})

			It("should record the admin by gorm set", func() {
				factory.Create(&paper, dbx.Opt{Set: h{"trail:who": "will"}})
				FindVersion()
				Expect(version.Whodunnit).To(Equal("will"))
			})
		})
		When("the paper is being updated", func() {
			It("creates the version which event is update, and computes the object changes", func() {
				factory.Create(&paper)
				factory.UpdateBy(&paper, Paper{Title: "abc", Other: "def"})
				FindVersion()
				Expect(version).To(HaveAttributes(Version{Event: "update", ItemID: paper.ID, ItemType: "Paper"}, "time"))
				Expect(version.Object).To(EqualYAML(h{"jsonb": "", "nilable": nil, "title": "abc"}))
				Expect(version.ObjectChanges).To(EqualYAML(h{"title": []interface{}{"test", "abc"}}))
			})

			It("should not create version if no changes", func() {
				factory.Create(&paper)
				factory.UpdateBy(&paper, h{"other": "def"})
				Expect(dbx.Model(&Version{}).Count()).To(BeEquivalentTo(1))
			})

			It("processes Jsonb", func() {
				factory.Create(&paper)
				factory.UpdateBy(&paper, h{"other": "def", "jsonb": structx.ToJsonb(h{"key": "value", "abc": "def"})})
				FindVersion()
				Expect(version.ObjectChanges).To(EqualYAML(h{"jsonb": []interface{}{"", "{\"abc\":\"def\",\"key\":\"value\"}"}}))

				factory.UpdateBy(&paper, h{"title": "xyz"})
				version = Version{}
				FindVersion()
				Expect(version.ObjectChanges).To(EqualYAML(h{"title": []interface{}{"test", "xyz"}}))

				Expect(dbx.Model(&Version{}).Count()).To(BeEquivalentTo(3))
				factory.UpdateBy(&paper, h{"jsonb": structx.ToJsonb(&h{"key": "value", "abc": "def"})})
				Expect(dbx.Model(&Version{}).Count()).To(BeEquivalentTo(3))

				Expect(paper.Trail(dbx.Conn().DB, &paper, h{"change": []interface{}{"a", "b"}}, "who")).To(Succeed())
				Expect(dbx.Model(&Version{}).Count()).To(BeEquivalentTo(4))
				factory.UpdateBy(&paper, h{"jsonb": structx.ToJsonb(&h{"key": "value", "abc": "def"})})
				Expect(dbx.Model(&Version{}).Count()).To(BeEquivalentTo(4))
			})

			It("processes nilable value", func() {
				factory.Create(&paper)
				factory.UpdateBy(&paper, h{"nilable": "nilable string"})
				FindVersion()
				Expect(version.ObjectChanges).To(EqualYAML(h{"nilable": []interface{}{nil, "nilable string"}}))
			})
		})
		When("the paper is being destroyed", func() {
			It("creates the version which event is destroy", func() {
				factory.Create(&paper)
				factory.Destroy(&paper)
				FindVersion()
				Expect(version).To(HaveAttributes(Version{Event: "destroy", ItemID: paper.ID, ItemType: "Paper"}, "time"))
			})
		})
	})

	Describe(".ChangeSet", func() {
		It("returns map", func() {
			factory.Create(&paper)
			FindVersion()
			Expect(version.ChangeSet()).To(Equal(map[string][]interface{}{"jsonb": {nil, ""}, "title": {nil, "test"}}))
		})
	})
})

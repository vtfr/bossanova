package service_test

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vtfr/bossanova/model"
	"github.com/vtfr/bossanova/service"
)

var _ = Describe("Authorization", func() {
	const samplePermissionData = `{
		"roles": {
			"roleA": [
				"permissionA",
				"permissionB"
			],
			"roleB": [
				"permissionC",
				"permissionD"
			],
			"default": [
				"permissionA"
			]
		},
		"inheritance": {
			"roleA": [ "roleB" ]
		}
	}`

	var auth service.Authorizator

	It("should fail on malformed configuration file", func() {
		r := strings.NewReader(`{
			"roles": {
				"roleA": true,
			}
		}`)

		_, err := service.NewAuthorizatorFromFile(r)
		Expect(err).To(HaveOccurred())
	})
	It("should read from the configuration file", func() {
		r := strings.NewReader(samplePermissionData)

		var err error
		auth, err = service.NewAuthorizatorFromFile(r)
		Expect(err).ToNot(HaveOccurred())
	})
	It("should authorize correctly", func() {
		userA := &model.User{Role: "roleA"}
		userB := &model.User{Role: "roleB"}

		Expect(auth.IsAuthorized(userA, "permissionA")).To(BeTrue())
		Expect(auth.IsAuthorized(userA, "permissionB")).To(BeTrue())
		Expect(auth.IsAuthorized(userA, "permissionC")).To(BeTrue())
		Expect(auth.IsAuthorized(userA, "permissionD")).To(BeTrue())

		Expect(auth.IsAuthorized(userB, "permissionA")).To(BeFalse())
		Expect(auth.IsAuthorized(userB, "permissionB")).To(BeFalse())
		Expect(auth.IsAuthorized(userB, "permissionC")).To(BeTrue())
		Expect(auth.IsAuthorized(userB, "permissionD")).To(BeTrue())
	})
	It("should authorize default correctly", func() {
		Expect(auth.IsAuthorized(nil, "permissionA")).To(BeTrue())
		Expect(auth.IsAuthorized(nil, "permissionB")).To(BeFalse())
		Expect(auth.IsAuthorized(nil, "permissionC")).To(BeFalse())
		Expect(auth.IsAuthorized(nil, "permissionD")).To(BeFalse())
	})
})

package validation_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	k6tv1 "kubevirt.io/client-go/api/v1"

	k6tobjs "github.com/fromanirh/kubevirt-template-validator/pkg/kubevirtobjs"
	"github.com/fromanirh/kubevirt-template-validator/pkg/validation"
)

var _ = Describe("Specialized", func() {
	Context("With valid data", func() {

		var (
			vmCirros *k6tv1.VirtualMachine
			vmRef    *k6tv1.VirtualMachine
		)

		BeforeEach(func() {
			vmCirros = NewVMCirros()
			vmRef = k6tobjs.NewDefaultVirtualMachine()
		})

		It("Should apply simple integer rules", func() {
			r := validation.Rule{
				Rule:    "integer",
				Name:    "EnoughMemory",
				Path:    "jsonpath::.spec.domain.resources.requests.memory",
				Message: "Memory size not specified",
				Min:     64 * 1024 * 1024,
				Max:     512 * 1024 * 1024,
			}
			expectRuleApplicationSuccess(&r, vmCirros, vmRef)
		})

		It("Should apply simple string rules", func() {
			r := validation.Rule{
				Rule:      "string",
				Name:      "HasChipset",
				Path:      "jsonpath::.spec.domain.machine.type",
				Message:   "machine type must be specified",
				MinLength: 1,
				MaxLength: 32,
			}
			expectRuleApplicationSuccess(&r, vmCirros, vmRef)
		})

		It("Should apply simple enum rules", func() {
			r := validation.Rule{
				Rule:    "enum",
				Name:    "SupportedChipset",
				Path:    "jsonpath::.spec.domain.machine.type",
				Message: "machine type must be a supported value",
				Values:  []string{"q35", "440fx"},
			}
			expectRuleApplicationSuccess(&r, vmCirros, vmRef)
		})

		It("Should apply enum rule to multiple values", func() {
			r := validation.Rule{
				Rule:    "enum",
				Name:    "SupportedDiskBus",
				Path:    "jsonpath::.spec.domain.devices.disks[*].disk.bus",
				Message: "disk bus must be a supported value",
				Values:  []string{"virtio", "sata"},
			}
			expectRuleApplicationSuccess(&r, vmCirros, vmRef)
		})

		It("Should apply simple regex rules", func() {
			r := validation.Rule{
				Rule:    "regex",
				Name:    "SupportedChipset",
				Path:    "jsonpath::.spec.domain.machine.type",
				Message: "machine type must be a supported value",
				Regex:   "q35|440fx",
			}
			expectRuleApplicationSuccess(&r, vmCirros, vmRef)
		})

		It("Should apply regex rule to multiple values", func() {
			r := validation.Rule{
				Rule:    "regex",
				Name:    "SupportedDiskBus",
				Path:    "jsonpath::.spec.domain.devices.disks[*].disk.bus",
				Message: "disk bus must be a supported value",
				Regex:   "virtio|sata",
			}
			expectRuleApplicationSuccess(&r, vmCirros, vmRef)
		})
	})

	Context("With invalid data", func() {

		var (
			vmCirros *k6tv1.VirtualMachine
			vmRef    *k6tv1.VirtualMachine
		)

		BeforeEach(func() {
			vmCirros = NewVMCirros()
			vmRef = k6tobjs.NewDefaultVirtualMachine()
		})

		It("Should detect bogus rules", func() {
			r := validation.Rule{
				Rule:    "integer-value",
				Name:    "EnoughMemory",
				Path:    "jsonpath::.spec.domain.resources.requests.memory",
				Message: "Memory size not specified",
				Valid:   "jsonpath::.spec.domain.this.path.does.not.exist",
				Min:     64 * 1024 * 1024,
				Max:     512 * 1024 * 1024,
			}

			ra, err := r.Specialize(vmCirros, vmRef)
			Expect(err).To(Not(BeNil()))
			Expect(ra).To(BeNil())
		})

		It("Should fail simple integer rules", func() {
			r1 := validation.Rule{
				Rule:    "integer",
				Name:    "EnoughMemory",
				Path:    "jsonpath::.spec.domain.resources.requests.memory",
				Message: "Memory size not specified",
				Valid:   "jsonpath::.spec.domain.this.path.does.not.exist",
				Min:     512 * 1024 * 1024,
			}
			expectRuleApplicationFailure(&r1, vmCirros, vmRef)

			r2 := validation.Rule{
				Rule:    "integer",
				Name:    "EnoughMemory",
				Path:    "jsonpath::.spec.domain.resources.requests.memory",
				Message: "Memory size not specified",
				Valid:   "jsonpath::.spec.domain.this.path.does.not.exist",
				Max:     64 * 1024 * 1024,
			}
			expectRuleApplicationFailure(&r2, vmCirros, vmRef)
		})

		It("Should apply simple string rules", func() {
			r1 := validation.Rule{
				Rule:      "string",
				Name:      "HasChipset",
				Path:      "jsonpath::.spec.domain.machine.type",
				Message:   "machine type must be specified",
				MinLength: 64,
			}
			expectRuleApplicationFailure(&r1, vmCirros, vmRef)

			r2 := validation.Rule{
				Rule:      "string",
				Name:      "HasChipset",
				Path:      "jsonpath::.spec.domain.machine.type",
				Message:   "machine type must be specified",
				MaxLength: 1,
			}
			expectRuleApplicationFailure(&r2, vmCirros, vmRef)
		})

		It("Should apply simple enum rules", func() {
			r := validation.Rule{
				Rule:    "enum",
				Name:    "SupportedChipset",
				Path:    "jsonpath::.spec.domain.machine.type",
				Message: "machine type must be a supported value",
				Values:  []string{"foo", "bar"},
			}
			expectRuleApplicationFailure(&r, vmCirros, vmRef)
		})

		It("Should apply enum rule to multiple values", func() {
			r := validation.Rule{
				Rule:    "enum",
				Name:    "SupportedDiskBus",
				Path:    "jsonpath::.spec.domain.devices.disks[*].disk.bus",
				Message: "disk bus must be a supported value",
				Values:  []string{"foo", "bar"},
			}
			expectRuleApplicationFailure(&r, vmCirros, vmRef)
		})

		It("Should error enum rule if values do not exist", func() {
			vmCirros.Spec.Template.Spec.Domain.Devices.Disks = nil
			r := validation.Rule{
				Rule:    "enum",
				Name:    "SupportedDiskBus",
				Path:    "jsonpath::.spec.domain.devices.disks[*].disk.bus",
				Message: "disk bus must be a supported value",
				Values:  []string{"virtio"},
			}
			expectRuleApplicationError(&r, vmCirros, vmRef)
		})

		It("Should apply simple regex rules", func() {
			r := validation.Rule{
				Rule:    "regex",
				Name:    "SupportedChipset",
				Path:    "jsonpath::.spec.domain.machine.type",
				Message: "machine type must be a supported value",
				Regex:   "\\d[a-z]+\\d\\d",
			}
			expectRuleApplicationFailure(&r, vmCirros, vmRef)
		})

		It("Should apply regex rule to multiple values", func() {
			r := validation.Rule{
				Rule:    "regex",
				Name:    "SupportedDiskBus",
				Path:    "jsonpath::.spec.domain.devices.disks[*].disk.bus",
				Message: "disk bus must be a supported value",
				Regex:   "foo|bar",
			}
			expectRuleApplicationFailure(&r, vmCirros, vmRef)
		})

		It("Should error regex rule if values do not exist", func() {
			vmCirros.Spec.Template.Spec.Domain.Devices.Disks = nil
			r := validation.Rule{
				Rule:    "regex",
				Name:    "SupportedDiskBus",
				Path:    "jsonpath::.spec.domain.devices.disks[*].disk.bus",
				Message: "disk bus must be a supported value",
				Regex:   "virtio|sata",
			}
			expectRuleApplicationError(&r, vmCirros, vmRef)
		})
	})

})

func expectRuleApplicationSuccess(r *validation.Rule, vm, ref *k6tv1.VirtualMachine) {
	checkRuleApplication(r, vm, ref, true)
}

func expectRuleApplicationFailure(r *validation.Rule, vm, ref *k6tv1.VirtualMachine) {
	checkRuleApplication(r, vm, ref, false)
}

func expectRuleApplicationError(r *validation.Rule, vm, ref *k6tv1.VirtualMachine) {
	ra, err := r.Specialize(vm, ref)
	Expect(err).To(BeNil())
	Expect(ra).To(Not(BeNil()))

	_, err = ra.Apply(vm, ref)
	Expect(err).To(HaveOccurred())
}

func checkRuleApplication(r *validation.Rule, vm, ref *k6tv1.VirtualMachine, expected bool) {
	ra, err := r.Specialize(vm, ref)
	Expect(err).To(BeNil())
	Expect(ra).To(Not(BeNil()))

	ok, err := ra.Apply(vm, ref)
	Expect(err).To(BeNil())
	Expect(ok).To(Equal(expected))
}

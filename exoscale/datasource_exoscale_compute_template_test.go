package exoscale

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const (
	dataSourceComputeTemplateID       = testInstanceTemplateID
	dataSourceComputeTemplateName     = "Linux Ubuntu 20.04 LTS 64-bit"
	dataSourceComputeTemplateUsername = "ubuntu"
)

func TestAccDataSourceComputeTemplate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
data "exoscale_compute_template" "ubuntu_lts" {
  zone = "ch-gva-2"
}`,
				ExpectError: regexp.MustCompile("either name or id must be specified"),
			},
			{
				Config: fmt.Sprintf(`
data "exoscale_compute_template" "ubuntu_lts" {
  zone   = "ch-gva-2"
  name   = "%s"
  filter = "featured"
}`, dataSourceComputeTemplateName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceComputeTemplateAttributes(testAttrs{
						"id":       ValidateString(dataSourceComputeTemplateID),
						"name":     ValidateString(dataSourceComputeTemplateName),
						"username": ValidateString(dataSourceComputeTemplateUsername),
					}),
				),
			},
			{
				Config: fmt.Sprintf(`
data "exoscale_compute_template" "ubuntu_lts" {
  zone   = "ch-gva-2"
  id     = "%s"
  filter = "featured"
}`, dataSourceComputeTemplateID),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceComputeTemplateAttributes(testAttrs{
						"id":       ValidateString(dataSourceComputeTemplateID),
						"name":     ValidateString(dataSourceComputeTemplateName),
						"username": ValidateString(dataSourceComputeTemplateUsername),
					}),
				),
			},
		},
	})
}

func testAccDataSourceComputeTemplateAttributes(expected testAttrs) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "exoscale_compute_template" {
				continue
			}

			return checkResourceAttributes(expected, rs.Primary.Attributes)
		}

		return errors.New("exoscale_compute_template data source not found in the state")
	}
}

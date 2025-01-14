package oncall_test

import (
	"fmt"
	"testing"

	onCallAPI "github.com/grafana/amixr-api-go-client"
	"github.com/grafana/terraform-provider-grafana/internal/common"
	"github.com/grafana/terraform-provider-grafana/internal/testutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOnCallIntegration_basic(t *testing.T) {
	testutils.CheckCloudInstanceTestsEnabled(t)

	rName := fmt.Sprintf("test-acc-%s", acctest.RandString(8))
	rType := "grafana"

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: testutils.ProviderFactories,
		CheckDestroy:      testAccCheckOnCallIntegrationResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOnCallIntegrationConfig(rName, rType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOnCallIntegrationResourceExists("grafana_oncall_integration.test-acc-integration"),
					resource.TestCheckResourceAttr("grafana_oncall_integration.test-acc-integration", "name", rName),
					resource.TestCheckResourceAttr("grafana_oncall_integration.test-acc-integration", "type", rType),
					resource.TestCheckResourceAttrSet("grafana_oncall_integration.test-acc-integration", "link"),
				),
			},
		},
	})
}

func testAccCheckOnCallIntegrationResourceDestroy(s *terraform.State) error {
	client := testutils.Provider.Meta().(*common.Client).OnCallClient
	for _, r := range s.RootModule().Resources {
		if r.Type != "grafana_oncall_integration" {
			continue
		}

		if _, _, err := client.Integrations.GetIntegration(r.Primary.ID, &onCallAPI.GetIntegrationOptions{}); err == nil {
			return fmt.Errorf("integration still exists")
		}
	}
	return nil
}

func testAccOnCallIntegrationConfig(rName, rType string) string {
	return fmt.Sprintf(`
resource "grafana_oncall_integration" "test-acc-integration" {
	name = "%s"
	type = "%s"
	default_route {
	    slack {
	        enabled = false
	    }
	    telegram {
	        enabled = false
	    }
	}
}
`, rName, rType)
}

func testAccCheckOnCallIntegrationResourceExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Integration ID is set")
		}

		client := testutils.Provider.Meta().(*common.Client).OnCallClient

		found, _, err := client.Integrations.GetIntegration(rs.Primary.ID, &onCallAPI.GetIntegrationOptions{})
		if err != nil {
			return err
		}
		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Integration not found: %v - %v", rs.Primary.ID, found)
		}
		return nil
	}
}

package evidently_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/cloudwatchevidently"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfcloudwatchevidently "github.com/hashicorp/terraform-provider-aws/internal/service/evidently"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccEvidentlyFeature_basic(t *testing.T) {
	var feature cloudwatchevidently.Feature

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rName2 := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_evidently_feature.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(cloudwatchevidently.EndpointsID, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, cloudwatchevidently.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFeatureDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFeatureConfig_basic(rName, rName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFeatureExists(resourceName, &feature),
					acctest.CheckResourceAttrRegionalARN(resourceName, "arn", "evidently", fmt.Sprintf("project/%s/feature/%s", rName, rName2)),
					resource.TestCheckResourceAttrSet(resourceName, "created_time"),
					resource.TestCheckResourceAttr(resourceName, "default_variation", "Variation1"),
					resource.TestCheckResourceAttr(resourceName, "entity_overrides.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "evaluation_rules.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "evaluation_strategy", cloudwatchevidently.FeatureEvaluationStrategyAllRules),
					resource.TestCheckResourceAttrSet(resourceName, "last_updated_time"),
					resource.TestCheckResourceAttr(resourceName, "name", rName2),
					resource.TestCheckResourceAttr(resourceName, "project", rName),
					resource.TestCheckResourceAttr(resourceName, "status", cloudwatchevidently.FeatureStatusAvailable),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "value_type", cloudwatchevidently.VariationValueTypeString),
					resource.TestCheckResourceAttr(resourceName, "variations.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "variations.*", map[string]string{
						"name":                 "Variation1",
						"value.#":              "1",
						"value.0.string_value": "test",
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccEvidentlyFeature_updateDefaultVariation(t *testing.T) {
	var feature cloudwatchevidently.Feature

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rName2 := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	variationName1 := "Variation1"
	variationName2 := "Variation2"
	resourceName := "aws_evidently_feature.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(cloudwatchevidently.EndpointsID, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, cloudwatchevidently.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFeatureDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFeatureConfig_defaultVariation(rName, rName2, variationName1, variationName2, "first"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFeatureExists(resourceName, &feature),
					resource.TestCheckResourceAttr(resourceName, "default_variation", variationName1),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccFeatureConfig_defaultVariation(rName, rName2, variationName1, variationName2, "second"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFeatureExists(resourceName, &feature),
					resource.TestCheckResourceAttr(resourceName, "default_variation", variationName2),
				),
			},
		},
	})
}

func TestAccEvidentlyFeature_updateDescription(t *testing.T) {
	var feature cloudwatchevidently.Feature

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rName2 := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	originalDescription := "original description"
	updatedDescription := "updated description"
	resourceName := "aws_evidently_feature.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(cloudwatchevidently.EndpointsID, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, cloudwatchevidently.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFeatureDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFeatureConfig_description(rName, rName2, originalDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFeatureExists(resourceName, &feature),
					resource.TestCheckResourceAttr(resourceName, "description", originalDescription),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccFeatureConfig_description(rName, rName2, updatedDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFeatureExists(resourceName, &feature),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
				),
			},
		},
	})
}

func TestAccEvidentlyFeature_updateEntityOverrides(t *testing.T) {
	var feature cloudwatchevidently.Feature

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rName2 := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	variationName1 := "Variation1"
	variationName2 := "Variation2"
	resourceName := "aws_evidently_feature.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(cloudwatchevidently.EndpointsID, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, cloudwatchevidently.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFeatureDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFeatureConfig_entityOverrides1(rName, rName2, variationName1, variationName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFeatureExists(resourceName, &feature),
					resource.TestCheckResourceAttr(resourceName, "entity_overrides.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "entity_overrides.test1", variationName1),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "variations.*", map[string]string{
						"name":                 variationName1,
						"value.#":              "1",
						"value.0.string_value": "testval1",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "variations.*", map[string]string{
						"name":                 variationName2,
						"value.#":              "1",
						"value.0.string_value": "testval2",
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccFeatureConfig_entityOverrides2(rName, rName2, variationName1, variationName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFeatureExists(resourceName, &feature),
					resource.TestCheckResourceAttr(resourceName, "entity_overrides.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "entity_overrides.test1", variationName2),
					resource.TestCheckResourceAttr(resourceName, "entity_overrides.test2", variationName1),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "variations.*", map[string]string{
						"name":                 variationName1,
						"value.#":              "1",
						"value.0.string_value": "testval1",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "variations.*", map[string]string{
						"name":                 variationName2,
						"value.#":              "1",
						"value.0.string_value": "testval2",
					}),
				),
			},
		},
	})
}

func TestAccEvidentlyFeature_updateEvaluationStrategy(t *testing.T) {
	var feature cloudwatchevidently.Feature

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rName2 := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	originalEvaluationStategy := cloudwatchevidently.FeatureEvaluationStrategyAllRules
	updatedEvaluationStategy := cloudwatchevidently.FeatureEvaluationStrategyDefaultVariation
	resourceName := "aws_evidently_feature.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(cloudwatchevidently.EndpointsID, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, cloudwatchevidently.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFeatureDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFeatureConfig_evaluationStrategy(rName, rName2, originalEvaluationStategy),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFeatureExists(resourceName, &feature),
					resource.TestCheckResourceAttr(resourceName, "evaluation_strategy", originalEvaluationStategy),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccFeatureConfig_evaluationStrategy(rName, rName2, updatedEvaluationStategy),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFeatureExists(resourceName, &feature),
					resource.TestCheckResourceAttr(resourceName, "evaluation_strategy", updatedEvaluationStategy),
				),
			},
		},
	})
}

func TestAccEvidentlyFeature_tags(t *testing.T) {
	var feature cloudwatchevidently.Feature

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rName2 := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_evidently_feature.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(cloudwatchevidently.EndpointsID, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, cloudwatchevidently.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFeatureDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFeatureConfig_tags1(rName, rName2, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFeatureExists(resourceName, &feature),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccFeatureConfig_tags2(rName, rName2, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFeatureExists(resourceName, &feature),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccFeatureConfig_tags1(rName, rName2, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFeatureExists(resourceName, &feature),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccEvidentlyFeature_disappears(t *testing.T) {
	var feature cloudwatchevidently.Feature

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rName2 := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_evidently_feature.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, cloudwatchevidently.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFeatureDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFeatureConfig_basic(rName, rName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFeatureExists(resourceName, &feature),
					acctest.CheckResourceDisappears(acctest.Provider, tfcloudwatchevidently.ResourceFeature(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckFeatureDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).EvidentlyConn
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_evidently_feature" {
			continue
		}

		featureName, projectNameOrARN, err := tfcloudwatchevidently.FeatureParseID(rs.Primary.ID)

		if err != nil {
			return err
		}

		_, err = tfcloudwatchevidently.FindFeatureWithProjectNameorARN(context.Background(), conn, featureName, projectNameOrARN)

		if tfresource.NotFound(err) {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("CloudWatch Evidently Feature %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccCheckFeatureExists(n string, v *cloudwatchevidently.Feature) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No CloudWatch Evidently Feature ID is set")
		}

		featureName, projectNameOrARN, err := tfcloudwatchevidently.FeatureParseID(rs.Primary.ID)

		if err != nil {
			return err
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).EvidentlyConn

		output, err := tfcloudwatchevidently.FindFeatureWithProjectNameorARN(context.Background(), conn, featureName, projectNameOrARN)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccFeatureConfigBase(rName string) string {
	return fmt.Sprintf(`
resource "aws_evidently_project" "test" {
  name = %[1]q
}
`, rName)
}

func testAccFeatureConfig_basic(rName, rName2 string) string {
	return acctest.ConfigCompose(
		testAccFeatureConfigBase(rName),
		fmt.Sprintf(`
resource "aws_evidently_feature" "test" {
  name    = %[1]q
  project = aws_evidently_project.test.name

  variations {
    name = "Variation1"
    value {
      string_value = "test"
    }
  }
}
`, rName2))
}

func testAccFeatureConfig_defaultVariation(rName, rName2, variationName1, variationName2, selectDefaultVariation string) string {
	return acctest.ConfigCompose(
		testAccFeatureConfigBase(rName),
		fmt.Sprintf(`
locals {
  select_default_variation = %[4]q
  variation_name1          = %[2]q
  variation_name2          = %[3]q
}

resource "aws_evidently_feature" "test" {
  name              = %[1]q
  project           = aws_evidently_project.test.name
  default_variation = local.select_default_variation == "first" ? local.variation_name1 : local.variation_name2

  variations {
    name = %[2]q
    value {
      string_value = "testval1"
    }
  }

  variations {
    name = %[3]q
    value {
      string_value = "testval2"
    }
  }
}
`, rName2, variationName1, variationName2, selectDefaultVariation))
}

func testAccFeatureConfig_description(rName, rName2, description string) string {
	return acctest.ConfigCompose(
		testAccFeatureConfigBase(rName),
		fmt.Sprintf(`
resource "aws_evidently_feature" "test" {
  name        = %[1]q
  description = %[2]q
  project     = aws_evidently_project.test.name

  variations {
    name = "Variation1"
    value {
      string_value = "test"
    }
  }
}
`, rName2, description))
}

func testAccFeatureConfig_entityOverrides1(rName, rName2, variationName1, variationName2 string) string {
	return acctest.ConfigCompose(
		testAccFeatureConfigBase(rName),
		fmt.Sprintf(`
resource "aws_evidently_feature" "test" {
  name    = %[1]q
  project = aws_evidently_project.test.name

  entity_overrides = {
    test1 = %[2]q
  }

  variations {
    name = %[2]q
    value {
      string_value = "testval1"
    }
  }

  variations {
    name = %[3]q
    value {
      string_value = "testval2"
    }
  }
}
`, rName2, variationName1, variationName2))
}

func testAccFeatureConfig_entityOverrides2(rName, rName2, variationName1, variationName2 string) string {
	return acctest.ConfigCompose(
		testAccFeatureConfigBase(rName),
		fmt.Sprintf(`
resource "aws_evidently_feature" "test" {
  name    = %[1]q
  project = aws_evidently_project.test.name

  entity_overrides = {
    test1 = %[3]q
    test2 = %[2]q
  }

  variations {
    name = %[2]q
    value {
      string_value = "testval1"
    }
  }

  variations {
    name = %[3]q
    value {
      string_value = "testval2"
    }
  }
}
`, rName2, variationName1, variationName2))
}

func testAccFeatureConfig_evaluationStrategy(rName, rName2, evaluationStrategy string) string {
	return acctest.ConfigCompose(
		testAccFeatureConfigBase(rName),
		fmt.Sprintf(`
resource "aws_evidently_feature" "test" {
  name                = %[1]q
  evaluation_strategy = %[2]q
  project             = aws_evidently_project.test.name

  variations {
    name = "Variation1"
    value {
      string_value = "test"
    }
  }
}
`, rName2, evaluationStrategy))
}

func testAccFeatureConfig_tags1(rName, rName2, tag, value string) string {
	return acctest.ConfigCompose(
		testAccFeatureConfigBase(rName),
		fmt.Sprintf(`
resource "aws_evidently_feature" "test" {
  name    = %[1]q
  project = aws_evidently_project.test.name

  variations {
    name = "Variation1"
    value {
      string_value = "test"
    }
  }

  tags = {
    %[2]q = %[3]q
  }
}
`, rName2, tag, value))
}

func testAccFeatureConfig_tags2(rName, rName2, tag1, value1, tag2, value2 string) string {
	return acctest.ConfigCompose(
		testAccFeatureConfigBase(rName),
		fmt.Sprintf(`
resource "aws_evidently_feature" "test" {
  name    = %[1]q
  project = aws_evidently_project.test.name

  variations {
    name = "Variation1"
    value {
      string_value = "test"
    }
  }

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName2, tag1, value1, tag2, value2))
}

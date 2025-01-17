package yandex

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/hashicorp/vault/helper/pgpkeys"
	"github.com/stretchr/testify/assert"

	"github.com/yandex-cloud/go-genproto/yandex/cloud/access"
)

func TestJoinedStrings(t *testing.T) {
	testKeys := []string{"key1", "key2", "key3"}
	joinedKeys := getJoinedKeys(testKeys)
	assert.Equal(t, "`key1`, `key2`, `key3`", joinedKeys)

	testKey := []string{"key1"}
	joinedKey := getJoinedKeys(testKey)
	assert.Equal(t, "`key1`", joinedKey)
}

func TestValidateParsableValue(t *testing.T) {
	correctParseFunc1 := func(value string) error {
		if value != "CORRECT" {
			return fmt.Errorf("expected correct value")
		}
		return nil
	}
	validator1 := validateParsableValue(correctParseFunc1)

	_, es := validator1("CORRECT", "some_key")
	assert.Equal(t, 0, len(es))

	_, es = validator1("INCORRECT", "some_key")
	assert.Equal(t, 1, len(es))

	_, es = validator1([]string{"wrong", "type", "should", "not", "panic"}, "some_key")
	assert.Equal(t, 1, len(es))

	_, es = validator1(666, "some_key")
	assert.Equal(t, 1, len(es))

	correctParseFunc2 := func(value int) (string, error) {
		if value < 500 {
			return "", fmt.Errorf("expected int >= 500")
		}
		return strconv.Itoa(value), nil
	}
	validator2 := validateParsableValue(correctParseFunc2)

	_, es = validator2(777, "some_key")
	assert.Equal(t, 0, len(es))

	_, es = validator2(99, "some_key")
	assert.Equal(t, 1, len(es))

	incorrectParseFunc := func() string {
		return "should not panic"
	}
	validator3 := validateParsableValue(incorrectParseFunc)

	_, es = validator3("something", "some_key")
	assert.Equal(t, 1, len(es))
}

func memberType(ab *access.AccessBinding) string {
	return ab.Subject.Type
}

func memberID(ab *access.AccessBinding) string {
	return ab.Subject.Id
}

func userAccountIDs(p *Policy) []string {
	usersMap := map[string]bool{}
	for _, b := range p.Bindings {
		if memberType(b) == "userAccount" {
			usersMap[memberID(b)] = true
		}
	}
	userIDs := []string{}
	for userID := range usersMap {
		userIDs = append(userIDs, userID)
	}
	return userIDs
}

func testAccCloudAssignCloudMemberRole(cloudID string, usersID ...string) (string, string) {
	var config bytes.Buffer
	var resourceRefs []string

	for _, userID := range usersID {
		resType := "yandex_resourcemanager_cloud_iam_member"
		resName := fmt.Sprintf("membership-%s-%s", cloudID, userID)

		config.WriteString(fmt.Sprintf(`
// Make user member of cloud to allow assign another roles
resource "%s" "%s" {
  cloud_id = "%s"
  role     = "resource-manager.clouds.member"
  member   = "userAccount:%s"
}
`, resType, resName, cloudID, userID))

		resourceRefs = append(resourceRefs, fmt.Sprintf("\"%s.%s\"", resType, resName))
	}

	return config.String(), strings.Join(resourceRefs, ",")
}

func getFolderIamPolicyByFolderID(folderID string, config *Config) (*Policy, error) {
	f := FolderIamUpdater{
		folderID: folderID,
		Config:   config,
	}
	return f.GetResourceIamPolicy()
}

func testAccCheckServiceAccountIam(resourceName, role string, members []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("can't find %s in state", resourceName)
		}

		bindings, err := getServiceAccountAccessBindings(config, rs.Primary.ID)
		if err != nil {
			return err
		}

		var roleMembers []string
		for _, binding := range bindings {
			if binding.RoleId == role {
				member := binding.Subject.Type + ":" + binding.Subject.Id
				roleMembers = append(roleMembers, member)
			}
		}
		sort.Strings(members)
		sort.Strings(roleMembers)

		if reflect.DeepEqual(members, roleMembers) {
			return nil
		}

		return fmt.Errorf("Binding found but expected members is %v, got %v", members, roleMembers)
	}
}

func testAccCheckCreatedAtAttr(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		const createdAtAttrName = "created_at"
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("can't find %s in state", resourceName)
		}

		createdAt, ok := rs.Primary.Attributes[createdAtAttrName]
		if !ok {
			return fmt.Errorf("can't find '%s' attr for %s resource", createdAtAttrName, resourceName)
		}

		if _, err := time.Parse(time.RFC3339, createdAt); err != nil {
			return fmt.Errorf("can't parse timestamp in attr '%s': %s", createdAtAttrName, createdAt)
		}
		return nil
	}
}

func testAccCheckResourceIDField(resourceName string, idFieldName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		if rs.Primary.Attributes[idFieldName] != rs.Primary.ID {
			return fmt.Errorf("Resource: %s id field: %s, doesn't match resource ID", resourceName, idFieldName)
		}

		return nil
	}
}

func testExistsElementWithAttrValue(resourceName, path, field, value string, fullPath *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ms := s.RootModule()
		rs, ok := ms.Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s in %s", resourceName, ms.Path)
		}

		is := rs.Primary
		if is == nil {
			return fmt.Errorf("No primary instance: %s in %s", resourceName, ms.Path)
		}

		for k, v := range is.Attributes {
			reStr := fmt.Sprintf(`(%s\.\d+)\.%s`, path, field)
			re := regexp.MustCompile(reStr)
			if re.MatchString(k) && v == value {
				sm := re.FindStringSubmatch(k)
				*fullPath = sm[1]
				return nil
			}
		}

		return fmt.Errorf(
			"Can't find key %s.*.%s in resource: %s with value %s", path, field, resourceName, value,
		)
	}
}

func testExistsFirstElementWithAttr(resourceName, path, field string, fullPath *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ms := s.RootModule()
		rs, ok := ms.Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s in %s", resourceName, ms.Path)
		}

		is := rs.Primary
		if is == nil {
			return fmt.Errorf("No primary instance: %s in %s", resourceName, ms.Path)
		}

		for k := range is.Attributes {
			reStr := fmt.Sprintf(`(%s\.\d+)\.%s`, path, field)
			re := regexp.MustCompile(reStr)
			if re.MatchString(k) {
				sm := re.FindStringSubmatch(k)
				*fullPath = sm[1]
				return nil
			}
		}

		return fmt.Errorf(
			"Can't find key %s.*.%s in resource: %s", path, field, resourceName,
		)
	}
}

func testCheckResourceSubAttr(resourceName string, path *string, field string, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ms := s.RootModule()
		rs, ok := ms.Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s in %s", resourceName, ms.Path)
		}

		is := rs.Primary
		if is == nil {
			return fmt.Errorf("No primary instance: %s in %s", resourceName, ms.Path)
		}

		fullPath := fmt.Sprintf("%s.%s", *path, field)
		actualValue, ok := is.Attributes[fullPath]
		if !ok {
			return fmt.Errorf("Can't find path %s in resource: %s", fullPath, resourceName)
		}

		if actualValue != value {
			return fmt.Errorf(
				"Can't match values for path %s in resource: %s. %s != %s", fullPath, resourceName, value, actualValue,
			)
		}

		return nil
	}
}

func testCheckResourceSubAttrFn(resourceName string, path *string, field string, checkfn func(string) error) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ms := s.RootModule()
		rs, ok := ms.Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s in %s", resourceName, ms.Path)
		}

		is := rs.Primary
		if is == nil {
			return fmt.Errorf("No primary instance: %s in %s", resourceName, ms.Path)
		}

		fullPath := fmt.Sprintf("%s.%s", *path, field)
		value, ok := is.Attributes[fullPath]
		if !ok {
			return fmt.Errorf("Can't find path %s in resource: %s", fullPath, resourceName)
		}

		err := checkfn(value)
		if err != nil {
			return err
		}

		return nil
	}
}

func testDecryptKeyAndTest(name, key, pgpKey string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}

		ciphertext, ok := rs.Primary.Attributes[key]
		if !ok {
			return fmt.Errorf("can't find '%s' attr for %s resource", key, name)
		}

		// We can't verify that the decrypted ciphertext is correct, because we don't
		// have it. We can verify that decrypting it does not error
		_, err := pgpkeys.DecryptBytes(ciphertext, pgpKey)
		if err != nil {
			return fmt.Errorf("error decrypting ciphertext: %s", err)
		}

		return nil
	}
}

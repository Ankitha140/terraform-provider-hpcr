// Copyright 2022 IBM Corp.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package example

import (
	_ "embed"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-provider-hpcr/validation"
)

//go:embed example1.tf
var ConfigExample1 string

//go:embed example2.tf
var ConfigExample2 string

func TestAccTgz(t *testing.T) {

	folder, _ := filepath.Abs("../samples/nginx-golang")

	t.Setenv("TF_VAR_FOLDER", folder)

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: ConfigExample1,
				Check:  resource.TestMatchOutput("result", validation.Base64Re),
			},
		},
	})
}

func TestAccTgzEncrypted(t *testing.T) {

	folder, _ := filepath.Abs("../samples/nginx-golang")

	t.Setenv("TF_VAR_FOLDER", folder)

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: ConfigExample2,
				Check:  resource.TestMatchOutput("result", validation.TokenRe),
			},
		},
	})
}

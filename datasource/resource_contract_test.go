// Copyright 2022 IBM Corp.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.package datasource

package datasource

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/terraform-provider-hpcr/common"
	D "github.com/terraform-provider-hpcr/data"
	"github.com/terraform-provider-hpcr/fp"
	E "github.com/terraform-provider-hpcr/fp/either"
	F "github.com/terraform-provider-hpcr/fp/function"
)

//go:embed samples/contract1.yaml
var Contract1 string

func TestContract(t *testing.T) {
	data := make(map[string]any)

	// prepare input data
	data[common.KeyContract] = Contract1
	data[common.KeyCert] = D.DefaultCertificate

	res := F.Pipe3(
		data,
		CreateResourceDataMock,
		resourceEncContract(&defaultContext),
		E.ToError[fp.ResourceData],
	)

	assert.NoError(t, res)

	fmt.Println(data[common.KeyRendered])
}

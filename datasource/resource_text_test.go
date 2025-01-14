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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/terraform-provider-hpcr/common"
	D "github.com/terraform-provider-hpcr/data"
	"github.com/terraform-provider-hpcr/encrypt"
	"github.com/terraform-provider-hpcr/fp"
	E "github.com/terraform-provider-hpcr/fp/either"
	F "github.com/terraform-provider-hpcr/fp/function"
	"github.com/terraform-provider-hpcr/validation"
)

var defaultContext = Context{
	encrypt.DefaultEncryption(),
	"test",
}

func TestUnencryptedText(t *testing.T) {
	data := make(map[string]any)

	// prepare input data
	data[common.KeyText] = "sample text"

	res := F.Pipe3(
		data,
		CreateResourceDataMock,
		resourceText(&defaultContext),
		E.ToError[fp.ResourceData],
	)

	assert.NoError(t, res)
	assert.Equal(t, data[common.KeyText], data[common.KeyRendered])
}

func TestEncryptedText(t *testing.T) {
	data := make(map[string]any)

	// prepare input data
	data[common.KeyText] = "sample text"
	data[common.KeyCert] = D.DefaultCertificate

	encText := resourceEncText(&defaultContext)

	res := F.Pipe4(
		data,
		CreateResourceDataMock,
		encText,
		E.Chain(encText),
		E.ToError[fp.ResourceData],
	)

	assert.NoError(t, res)
	assert.Regexp(t, validation.TokenRe, data[common.KeyRendered])
}

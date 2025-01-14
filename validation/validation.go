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
package validation

import (
	"fmt"
	"io/fs"
	"os"
	"regexp"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/terraform-provider-hpcr/common"
	"github.com/terraform-provider-hpcr/encrypt"
	"github.com/terraform-provider-hpcr/fp"
	E "github.com/terraform-provider-hpcr/fp/either"
	F "github.com/terraform-provider-hpcr/fp/function"
)

var (
	statE    = E.Eitherize1(os.Stat)
	Base64Re = regexp.MustCompile(`^((?:[A-Za-z\d+/]{4})*(?:[A-Za-z\d+/]{3}=|[A-Za-z\d+/]{2}==)?)$`)
	TokenRe  = regexp.MustCompile(`^hyper-protect-basic\.((?:[A-Za-z\d+/]{4})*(?:[A-Za-z\d+/]{3}=|[A-Za-z\d+/]{2}==)?)\.((?:[A-Za-z\d+/]{4})*(?:[A-Za-z\d+/]{3}=|[A-Za-z\d+/]{2}==)?)$`)
)

func toDiagnostics[A any](value E.Either[error, A]) diag.Diagnostics {
	return F.Pipe1(
		value,
		E.Fold(diag.FromErr, F.Constant1[A, diag.Diagnostics](nil)),
	)
}

// DiagCertificate validates that the given certificate is indeed a certificate
func DiagCertificate(data any, _ cty.Path) diag.Diagnostics {
	// convert the key
	return F.Pipe4(
		data,
		fp.ToTypeE[string],
		common.MapStgToBytesE,
		E.Chain(encrypt.CertSerial),
		toDiagnostics[[]byte],
	)
}

// DiagFolder validates that the given path points to an existing folder
func DiagFolder(data any, _ cty.Path) diag.Diagnostics {
	return F.Pipe4(
		data,
		fp.ToTypeE[string],
		E.Chain(statE),
		E.Chain(E.FromPredicate(fs.FileInfo.IsDir, func(info fs.FileInfo) error {
			return fmt.Errorf("path %s is not a folder", info.Name())
		})),
		toDiagnostics[fs.FileInfo],
	)
}

//
// Licensed Materials - Property of IBM
//
// 5737-I09
//
// Copyright IBM Corp. 2022 All Rights Reserved.
// US Government Users Restricted Rights - Use, duplication or
// disclosure restricted by GSA ADP Schedule Contract with IBM Corp
//
package either

import (
	F "github.com/terraform-provider-hpcr/fp/function"
)

// constructs a function that creates a resource, then operates on it and then releases the resource
func WithResource[E, R, A any](onCreate func() Either[E, R], onRelease func(R) Either[E, any]) func(func(R) Either[E, A]) Either[E, A] {

	return func(f func(R) Either[E, A]) Either[E, A] {
		return MonadChain(
			onCreate(), func(r R) Either[E, A] {
				// run the code and make sure to release as quickly as possible
				res := f(r)
				released := onRelease(r)
				// handle the errors
				return fold(
					res,
					Left[E, A],
					func(a A) Either[E, A] {
						return F.Pipe1(
							released,
							MapTo[E, any](a),
						)
					})
			},
		)
	}
}
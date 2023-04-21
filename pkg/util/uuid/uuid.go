/*
Copyright 2022-2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package uuid

import (
	"github.com/google/uuid"
)

// Generate a new UUID, optionally with a prefix
func Generate(prefix string) string {
	id := uuid.New()
	if prefix != "" {
		return prefix + "-" + id.String()

	}
	return id.String()
}

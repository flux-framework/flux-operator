/*
Copyright 2022 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package controllers

import _ "embed"

//go:embed templates/broker.toml
var brokerConfigTemplate string

//go:embed templates/wait.sh
var waitToStartTemplate string

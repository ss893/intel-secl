/*
 * Copyright (C) 2020 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package hosttrust

import (
	"github.com/google/uuid"
	"github.com/intel-secl/intel-secl/v4/pkg/model/hvs"
)

type hostTrustCache struct {
	hostID         uuid.UUID
	trustedFlavors []hvs.Flavor
	// TODO: consider using a map here rather than traversing through the list when we need to remove flavors
	trustReport hvs.TrustReport
}

// TODO:
// These functions should have been implemented in
// "github.com/intel-secl/intel-secl/v4/pkg/model/hvs"
// for the structure hvs.FlavorCollection
func (htc *hostTrustCache) addTrustedFlavors(f *hvs.Flavor) {
	htc.trustedFlavors = append(htc.trustedFlavors, *f)
}

func (htc *hostTrustCache) removeTrustedFlavors(fIn *hvs.Flavor) {
	if fIn == nil {
		return
	}
	targetID := fIn.Meta.ID
	for i, f := range htc.trustedFlavors {
		if f.Meta.ID == targetID {
			htc.trustedFlavors = append(htc.trustedFlavors[:i], htc.trustedFlavors[i+1:]...)
		}
	}
}

func (htc *hostTrustCache) isTrustCacheEmpty() bool {
	if len(htc.trustedFlavors) == 0 {
		return true
	}
	return false
}

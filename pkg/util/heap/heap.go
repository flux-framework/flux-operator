/*
Copyright 2022 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package heap

import (
	jobctrl "flux-framework/flux-operator/pkg/job"
)

// lessFunc is a function that receives two items and returns true if the first
// item should be placed before the second one when the list is sorted.
type lessFunc func(a, b interface{}) bool

// KeyFunc is a function type to get the key from an object.
type keyFunc func(obj interface{}) string

type heapItem struct {
	obj   interface{}
	index int
	key   string
}

// A wrapper around a priority queue with
type Heap struct {
	// items is a map from key of the objects to the objects and their index
	items map[string]heapItem

	// Look up indices based on key
	keys     []string
	keyFunc  keyFunc
	lessFunc lessFunc
}

// Push a new job to the heap if not present
func (h *Heap) PushIfNotPresent(info *jobctrl.Info) bool {

	key := info.Obj.Name

	// If the JobInfo name isn't in items, add it
	if _, exists := h.items[key]; !exists {
		newItem := heapItem{info.Obj, len(h.items), key}
		h.items[key] = newItem
		return true
	}
	return false
}

// Push a new job to the heap if not present
func (h *Heap) Delete(info *jobctrl.Info) bool {

	key := info.Obj.Name

	// If the JobInfo name isn't in items, add it
	if _, exists := h.items[key]; exists {
		delete(h.items, key)
		return true
	}
	return false
}

// GetByKey gets a job based on the key
func (h *Heap) Exists(info *jobctrl.Info) bool {

	key := info.Obj.Name

	// If the JobInfo name isn't in items, add it
	_, exists := h.items[key]
	return exists
}

// Push or update
func (h *Heap) PushOrUpdate(info *jobctrl.Info) {

	key := info.Obj.Name

	// We add a new item no matter what
	newItem := heapItem{info.Obj, len(h.items), key}
	h.items[key] = newItem
}

func (h *Heap) Len() int {
	return len(h.items)
}

// New returns a Fake Heap to keep items
// This can eventually be a real heap
func New(keyFn keyFunc, lessFn lessFunc) Heap {
	items := map[string]heapItem{}
	return Heap{
		items:    items,
		keys:     []string{},
		keyFunc:  keyFn,
		lessFunc: lessFn,
	}
}

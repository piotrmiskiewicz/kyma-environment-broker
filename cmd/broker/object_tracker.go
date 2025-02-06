package main

import (
	"sync"
	"time"

	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	k8stesting "k8s.io/client-go/testing"
)

func NewTestingObjectTracker(sch *runtime.Scheme) *testingObjectTracker {
	return &testingObjectTracker{
		target:           k8stesting.NewObjectTracker(sch, scheme.Codecs.UniversalDecoder()),
		runtimesToDelete: map[string]struct{}{},
	}
}

type testingObjectTracker struct {
	target k8stesting.ObjectTracker

	mu               sync.RWMutex
	runtimesToDelete map[string]struct{}
}

type Deleter interface {
	ProcessRuntimeDeletion(name string)
}

var _ k8stesting.ObjectTracker = &testingObjectTracker{}
var _ Deleter = &testingObjectTracker{}

func (m *testingObjectTracker) Add(obj runtime.Object) error {
	return m.target.Add(obj)
}

func (m *testingObjectTracker) Get(gvr schema.GroupVersionResource, ns, name string) (runtime.Object, error) {
	return m.target.Get(gvr, ns, name)
}

func (m *testingObjectTracker) Create(gvr schema.GroupVersionResource, obj runtime.Object, ns string) error {
	return m.target.Create(gvr, obj, ns)
}

func (m *testingObjectTracker) Update(gvr schema.GroupVersionResource, obj runtime.Object, ns string) error {
	return m.target.Update(gvr, obj, ns)
}

func (m *testingObjectTracker) List(gvr schema.GroupVersionResource, gvk schema.GroupVersionKind, ns string) (runtime.Object, error) {
	return m.target.List(gvr, gvk, ns)
}

func (m *testingObjectTracker) Delete(gvr schema.GroupVersionResource, ns, name string) error {
	if gvr.Resource == "runtimes" {
		m.mu.Lock()
		defer m.mu.Unlock()
		m.runtimesToDelete[name] = struct{}{}
		return nil
	}
	return m.target.Delete(gvr, ns, name)
}

func (m *testingObjectTracker) Watch(gvr schema.GroupVersionResource, ns string) (watch.Interface, error) {
	return m.target.Watch(gvr, ns)
}

func (m *testingObjectTracker) ProcessRuntimeDeletion(name string) {
	for {
		time.Sleep(time.Millisecond)
		m.mu.RLock()
		_, ok := m.runtimesToDelete[name]
		m.mu.RUnlock()
		if ok {
			break
		}
	}
	m.deleteRuntimeIfExist(name)
}

func (m *testingObjectTracker) deleteRuntimeIfExist(name string) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if _, ok := m.runtimesToDelete[name]; ok {
		_ = m.target.Delete(imv1.GroupVersion.WithResource("runtimes"), "kyma-system", name)
	}
}

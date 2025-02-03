package main

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	k8stesting "k8s.io/client-go/testing"
	"sync"
	"time"

	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
)

func NewDeletionObjectTracker(sch *runtime.Scheme) *myOT {
	return &myOT{
		target:           k8stesting.NewObjectTracker(sch, scheme.Codecs.UniversalDecoder()),
		runtimesToDelete: map[string]struct{}{},
	}
}

type myOT struct {
	target k8stesting.ObjectTracker

	mu               sync.RWMutex
	runtimesToDelete map[string]struct{}
}

type RealDeleter interface {
	ProcessRuntimeDeletion(name string) error
}

var _ k8stesting.ObjectTracker = &myOT{}
var _ RealDeleter = &myOT{}

func (m *myOT) Add(obj runtime.Object) error {
	return m.target.Add(obj)
}

func (m *myOT) Get(gvr schema.GroupVersionResource, ns, name string) (runtime.Object, error) {
	return m.target.Get(gvr, ns, name)
}

func (m *myOT) Create(gvr schema.GroupVersionResource, obj runtime.Object, ns string) error {
	return m.target.Create(gvr, obj, ns)
}

func (m *myOT) Update(gvr schema.GroupVersionResource, obj runtime.Object, ns string) error {
	return m.target.Update(gvr, obj, ns)
}

func (m *myOT) List(gvr schema.GroupVersionResource, gvk schema.GroupVersionKind, ns string) (runtime.Object, error) {
	return m.target.List(gvr, gvk, ns)
}

func (m *myOT) Delete(gvr schema.GroupVersionResource, ns, name string) error {
	if gvr.Resource == "runtimes" {
		m.mu.Lock()
		defer m.mu.Unlock()
		m.runtimesToDelete[name] = struct{}{}
		return nil
	}
	return m.target.Delete(gvr, ns, name)
}

func (m *myOT) Watch(gvr schema.GroupVersionResource, ns string) (watch.Interface, error) {
	return m.target.Watch(gvr, ns)
}

func (m *myOT) ProcessRuntimeDeletion(name string) error {
	go func() {
		for {
			m.deleteRuntimeIfExeist(name)
			time.Sleep(1 * time.Millisecond)
		}
	}()
	return nil
}

func (m *myOT) deleteRuntimeIfExeist(name string) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if _, ok := m.runtimesToDelete[name]; ok {
		m.target.Delete(imv1.GroupVersion.WithResource("runtimes"), "kyma-system", name)
	}
}

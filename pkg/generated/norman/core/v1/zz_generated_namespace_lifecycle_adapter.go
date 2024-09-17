package v1

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type namespaceLifecycleConverter struct {
	lifecycle NamespaceLifecycle
}

func (w *namespaceLifecycleConverter) CreateContext(_ context.Context, obj *v1.Namespace) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *namespaceLifecycleConverter) RemoveContext(_ context.Context, obj *v1.Namespace) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *namespaceLifecycleConverter) UpdatedContext(_ context.Context, obj *v1.Namespace) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type NamespaceLifecycle interface {
	Create(obj *v1.Namespace) (runtime.Object, error)
	Remove(obj *v1.Namespace) (runtime.Object, error)
	Updated(obj *v1.Namespace) (runtime.Object, error)
}

type NamespaceLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v1.Namespace) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v1.Namespace) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v1.Namespace) (runtime.Object, error)
}

type namespaceLifecycleAdapter struct {
	lifecycle NamespaceLifecycleContext
}

func (w *namespaceLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *namespaceLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *namespaceLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *namespaceLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v1.Namespace))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *namespaceLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *namespaceLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v1.Namespace))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *namespaceLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *namespaceLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v1.Namespace))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewNamespaceLifecycleAdapter(name string, clusterScoped bool, client NamespaceInterface, l NamespaceLifecycle) NamespaceHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(NamespaceGroupVersionResource)
	}
	adapter := &namespaceLifecycleAdapter{lifecycle: &namespaceLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v1.Namespace) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewNamespaceLifecycleAdapterContext(name string, clusterScoped bool, client NamespaceInterface, l NamespaceLifecycleContext) NamespaceHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(NamespaceGroupVersionResource)
	}
	adapter := &namespaceLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v1.Namespace) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

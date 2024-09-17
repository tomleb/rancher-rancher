package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/project.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type namespacedServiceAccountTokenLifecycleConverter struct {
	lifecycle NamespacedServiceAccountTokenLifecycle
}

func (w *namespacedServiceAccountTokenLifecycleConverter) CreateContext(_ context.Context, obj *v3.NamespacedServiceAccountToken) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *namespacedServiceAccountTokenLifecycleConverter) RemoveContext(_ context.Context, obj *v3.NamespacedServiceAccountToken) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *namespacedServiceAccountTokenLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.NamespacedServiceAccountToken) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type NamespacedServiceAccountTokenLifecycle interface {
	Create(obj *v3.NamespacedServiceAccountToken) (runtime.Object, error)
	Remove(obj *v3.NamespacedServiceAccountToken) (runtime.Object, error)
	Updated(obj *v3.NamespacedServiceAccountToken) (runtime.Object, error)
}

type NamespacedServiceAccountTokenLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.NamespacedServiceAccountToken) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.NamespacedServiceAccountToken) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.NamespacedServiceAccountToken) (runtime.Object, error)
}

type namespacedServiceAccountTokenLifecycleAdapter struct {
	lifecycle NamespacedServiceAccountTokenLifecycleContext
}

func (w *namespacedServiceAccountTokenLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *namespacedServiceAccountTokenLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *namespacedServiceAccountTokenLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *namespacedServiceAccountTokenLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.NamespacedServiceAccountToken))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *namespacedServiceAccountTokenLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *namespacedServiceAccountTokenLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.NamespacedServiceAccountToken))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *namespacedServiceAccountTokenLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *namespacedServiceAccountTokenLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.NamespacedServiceAccountToken))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewNamespacedServiceAccountTokenLifecycleAdapter(name string, clusterScoped bool, client NamespacedServiceAccountTokenInterface, l NamespacedServiceAccountTokenLifecycle) NamespacedServiceAccountTokenHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(NamespacedServiceAccountTokenGroupVersionResource)
	}
	adapter := &namespacedServiceAccountTokenLifecycleAdapter{lifecycle: &namespacedServiceAccountTokenLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.NamespacedServiceAccountToken) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewNamespacedServiceAccountTokenLifecycleAdapterContext(name string, clusterScoped bool, client NamespacedServiceAccountTokenInterface, l NamespacedServiceAccountTokenLifecycleContext) NamespacedServiceAccountTokenHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(NamespacedServiceAccountTokenGroupVersionResource)
	}
	adapter := &namespacedServiceAccountTokenLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.NamespacedServiceAccountToken) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

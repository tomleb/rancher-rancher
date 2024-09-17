package v1

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type configMapLifecycleConverter struct {
	lifecycle ConfigMapLifecycle
}

func (w *configMapLifecycleConverter) CreateContext(_ context.Context, obj *v1.ConfigMap) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *configMapLifecycleConverter) RemoveContext(_ context.Context, obj *v1.ConfigMap) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *configMapLifecycleConverter) UpdatedContext(_ context.Context, obj *v1.ConfigMap) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type ConfigMapLifecycle interface {
	Create(obj *v1.ConfigMap) (runtime.Object, error)
	Remove(obj *v1.ConfigMap) (runtime.Object, error)
	Updated(obj *v1.ConfigMap) (runtime.Object, error)
}

type ConfigMapLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v1.ConfigMap) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v1.ConfigMap) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v1.ConfigMap) (runtime.Object, error)
}

type configMapLifecycleAdapter struct {
	lifecycle ConfigMapLifecycleContext
}

func (w *configMapLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *configMapLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *configMapLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *configMapLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v1.ConfigMap))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *configMapLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *configMapLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v1.ConfigMap))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *configMapLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *configMapLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v1.ConfigMap))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewConfigMapLifecycleAdapter(name string, clusterScoped bool, client ConfigMapInterface, l ConfigMapLifecycle) ConfigMapHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(ConfigMapGroupVersionResource)
	}
	adapter := &configMapLifecycleAdapter{lifecycle: &configMapLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v1.ConfigMap) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewConfigMapLifecycleAdapterContext(name string, clusterScoped bool, client ConfigMapInterface, l ConfigMapLifecycleContext) ConfigMapHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(ConfigMapGroupVersionResource)
	}
	adapter := &configMapLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v1.ConfigMap) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

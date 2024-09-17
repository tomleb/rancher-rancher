package v1

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type secretLifecycleConverter struct {
	lifecycle SecretLifecycle
}

func (w *secretLifecycleConverter) CreateContext(_ context.Context, obj *v1.Secret) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *secretLifecycleConverter) RemoveContext(_ context.Context, obj *v1.Secret) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *secretLifecycleConverter) UpdatedContext(_ context.Context, obj *v1.Secret) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type SecretLifecycle interface {
	Create(obj *v1.Secret) (runtime.Object, error)
	Remove(obj *v1.Secret) (runtime.Object, error)
	Updated(obj *v1.Secret) (runtime.Object, error)
}

type SecretLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v1.Secret) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v1.Secret) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v1.Secret) (runtime.Object, error)
}

type secretLifecycleAdapter struct {
	lifecycle SecretLifecycleContext
}

func (w *secretLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *secretLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *secretLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *secretLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v1.Secret))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *secretLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *secretLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v1.Secret))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *secretLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *secretLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v1.Secret))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewSecretLifecycleAdapter(name string, clusterScoped bool, client SecretInterface, l SecretLifecycle) SecretHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(SecretGroupVersionResource)
	}
	adapter := &secretLifecycleAdapter{lifecycle: &secretLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v1.Secret) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewSecretLifecycleAdapterContext(name string, clusterScoped bool, client SecretInterface, l SecretLifecycleContext) SecretHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(SecretGroupVersionResource)
	}
	adapter := &secretLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v1.Secret) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

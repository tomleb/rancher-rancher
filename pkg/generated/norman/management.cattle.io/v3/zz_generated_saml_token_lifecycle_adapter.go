package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type samlTokenLifecycleConverter struct {
	lifecycle SamlTokenLifecycle
}

func (w *samlTokenLifecycleConverter) CreateContext(_ context.Context, obj *v3.SamlToken) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *samlTokenLifecycleConverter) RemoveContext(_ context.Context, obj *v3.SamlToken) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *samlTokenLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.SamlToken) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type SamlTokenLifecycle interface {
	Create(obj *v3.SamlToken) (runtime.Object, error)
	Remove(obj *v3.SamlToken) (runtime.Object, error)
	Updated(obj *v3.SamlToken) (runtime.Object, error)
}

type SamlTokenLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.SamlToken) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.SamlToken) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.SamlToken) (runtime.Object, error)
}

type samlTokenLifecycleAdapter struct {
	lifecycle SamlTokenLifecycleContext
}

func (w *samlTokenLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *samlTokenLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *samlTokenLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *samlTokenLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.SamlToken))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *samlTokenLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *samlTokenLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.SamlToken))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *samlTokenLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *samlTokenLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.SamlToken))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewSamlTokenLifecycleAdapter(name string, clusterScoped bool, client SamlTokenInterface, l SamlTokenLifecycle) SamlTokenHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(SamlTokenGroupVersionResource)
	}
	adapter := &samlTokenLifecycleAdapter{lifecycle: &samlTokenLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.SamlToken) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewSamlTokenLifecycleAdapterContext(name string, clusterScoped bool, client SamlTokenInterface, l SamlTokenLifecycleContext) SamlTokenHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(SamlTokenGroupVersionResource)
	}
	adapter := &samlTokenLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.SamlToken) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

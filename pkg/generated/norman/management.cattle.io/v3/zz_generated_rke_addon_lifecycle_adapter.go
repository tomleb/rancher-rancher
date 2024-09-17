package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type rkeAddonLifecycleConverter struct {
	lifecycle RkeAddonLifecycle
}

func (w *rkeAddonLifecycleConverter) CreateContext(_ context.Context, obj *v3.RkeAddon) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *rkeAddonLifecycleConverter) RemoveContext(_ context.Context, obj *v3.RkeAddon) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *rkeAddonLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.RkeAddon) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type RkeAddonLifecycle interface {
	Create(obj *v3.RkeAddon) (runtime.Object, error)
	Remove(obj *v3.RkeAddon) (runtime.Object, error)
	Updated(obj *v3.RkeAddon) (runtime.Object, error)
}

type RkeAddonLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.RkeAddon) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.RkeAddon) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.RkeAddon) (runtime.Object, error)
}

type rkeAddonLifecycleAdapter struct {
	lifecycle RkeAddonLifecycleContext
}

func (w *rkeAddonLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *rkeAddonLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *rkeAddonLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *rkeAddonLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.RkeAddon))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *rkeAddonLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *rkeAddonLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.RkeAddon))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *rkeAddonLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *rkeAddonLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.RkeAddon))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewRkeAddonLifecycleAdapter(name string, clusterScoped bool, client RkeAddonInterface, l RkeAddonLifecycle) RkeAddonHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(RkeAddonGroupVersionResource)
	}
	adapter := &rkeAddonLifecycleAdapter{lifecycle: &rkeAddonLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.RkeAddon) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewRkeAddonLifecycleAdapterContext(name string, clusterScoped bool, client RkeAddonInterface, l RkeAddonLifecycleContext) RkeAddonHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(RkeAddonGroupVersionResource)
	}
	adapter := &rkeAddonLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.RkeAddon) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

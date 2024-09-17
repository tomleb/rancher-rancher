package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type globalDnsLifecycleConverter struct {
	lifecycle GlobalDnsLifecycle
}

func (w *globalDnsLifecycleConverter) CreateContext(_ context.Context, obj *v3.GlobalDns) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *globalDnsLifecycleConverter) RemoveContext(_ context.Context, obj *v3.GlobalDns) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *globalDnsLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.GlobalDns) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type GlobalDnsLifecycle interface {
	Create(obj *v3.GlobalDns) (runtime.Object, error)
	Remove(obj *v3.GlobalDns) (runtime.Object, error)
	Updated(obj *v3.GlobalDns) (runtime.Object, error)
}

type GlobalDnsLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.GlobalDns) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.GlobalDns) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.GlobalDns) (runtime.Object, error)
}

type globalDnsLifecycleAdapter struct {
	lifecycle GlobalDnsLifecycleContext
}

func (w *globalDnsLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *globalDnsLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *globalDnsLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *globalDnsLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.GlobalDns))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *globalDnsLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *globalDnsLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.GlobalDns))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *globalDnsLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *globalDnsLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.GlobalDns))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewGlobalDnsLifecycleAdapter(name string, clusterScoped bool, client GlobalDnsInterface, l GlobalDnsLifecycle) GlobalDnsHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(GlobalDnsGroupVersionResource)
	}
	adapter := &globalDnsLifecycleAdapter{lifecycle: &globalDnsLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.GlobalDns) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewGlobalDnsLifecycleAdapterContext(name string, clusterScoped bool, client GlobalDnsInterface, l GlobalDnsLifecycleContext) GlobalDnsHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(GlobalDnsGroupVersionResource)
	}
	adapter := &globalDnsLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.GlobalDns) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

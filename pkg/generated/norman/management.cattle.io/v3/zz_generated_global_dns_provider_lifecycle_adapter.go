package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type globalDnsProviderLifecycleConverter struct {
	lifecycle GlobalDnsProviderLifecycle
}

func (w *globalDnsProviderLifecycleConverter) CreateContext(_ context.Context, obj *v3.GlobalDnsProvider) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *globalDnsProviderLifecycleConverter) RemoveContext(_ context.Context, obj *v3.GlobalDnsProvider) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *globalDnsProviderLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.GlobalDnsProvider) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type GlobalDnsProviderLifecycle interface {
	Create(obj *v3.GlobalDnsProvider) (runtime.Object, error)
	Remove(obj *v3.GlobalDnsProvider) (runtime.Object, error)
	Updated(obj *v3.GlobalDnsProvider) (runtime.Object, error)
}

type GlobalDnsProviderLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.GlobalDnsProvider) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.GlobalDnsProvider) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.GlobalDnsProvider) (runtime.Object, error)
}

type globalDnsProviderLifecycleAdapter struct {
	lifecycle GlobalDnsProviderLifecycleContext
}

func (w *globalDnsProviderLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *globalDnsProviderLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *globalDnsProviderLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *globalDnsProviderLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.GlobalDnsProvider))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *globalDnsProviderLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *globalDnsProviderLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.GlobalDnsProvider))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *globalDnsProviderLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *globalDnsProviderLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.GlobalDnsProvider))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewGlobalDnsProviderLifecycleAdapter(name string, clusterScoped bool, client GlobalDnsProviderInterface, l GlobalDnsProviderLifecycle) GlobalDnsProviderHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(GlobalDnsProviderGroupVersionResource)
	}
	adapter := &globalDnsProviderLifecycleAdapter{lifecycle: &globalDnsProviderLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.GlobalDnsProvider) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewGlobalDnsProviderLifecycleAdapterContext(name string, clusterScoped bool, client GlobalDnsProviderInterface, l GlobalDnsProviderLifecycleContext) GlobalDnsProviderHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(GlobalDnsProviderGroupVersionResource)
	}
	adapter := &globalDnsProviderLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.GlobalDnsProvider) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

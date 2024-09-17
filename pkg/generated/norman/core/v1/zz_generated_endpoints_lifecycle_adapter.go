package v1

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type endpointsLifecycleConverter struct {
	lifecycle EndpointsLifecycle
}

func (w *endpointsLifecycleConverter) CreateContext(_ context.Context, obj *v1.Endpoints) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *endpointsLifecycleConverter) RemoveContext(_ context.Context, obj *v1.Endpoints) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *endpointsLifecycleConverter) UpdatedContext(_ context.Context, obj *v1.Endpoints) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type EndpointsLifecycle interface {
	Create(obj *v1.Endpoints) (runtime.Object, error)
	Remove(obj *v1.Endpoints) (runtime.Object, error)
	Updated(obj *v1.Endpoints) (runtime.Object, error)
}

type EndpointsLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v1.Endpoints) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v1.Endpoints) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v1.Endpoints) (runtime.Object, error)
}

type endpointsLifecycleAdapter struct {
	lifecycle EndpointsLifecycleContext
}

func (w *endpointsLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *endpointsLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *endpointsLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *endpointsLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v1.Endpoints))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *endpointsLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *endpointsLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v1.Endpoints))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *endpointsLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *endpointsLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v1.Endpoints))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewEndpointsLifecycleAdapter(name string, clusterScoped bool, client EndpointsInterface, l EndpointsLifecycle) EndpointsHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(EndpointsGroupVersionResource)
	}
	adapter := &endpointsLifecycleAdapter{lifecycle: &endpointsLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v1.Endpoints) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewEndpointsLifecycleAdapterContext(name string, clusterScoped bool, client EndpointsInterface, l EndpointsLifecycleContext) EndpointsHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(EndpointsGroupVersionResource)
	}
	adapter := &endpointsLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v1.Endpoints) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

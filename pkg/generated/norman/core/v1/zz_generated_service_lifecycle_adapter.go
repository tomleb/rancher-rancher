package v1

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type serviceLifecycleConverter struct {
	lifecycle ServiceLifecycle
}

func (w *serviceLifecycleConverter) CreateContext(_ context.Context, obj *v1.Service) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *serviceLifecycleConverter) RemoveContext(_ context.Context, obj *v1.Service) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *serviceLifecycleConverter) UpdatedContext(_ context.Context, obj *v1.Service) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type ServiceLifecycle interface {
	Create(obj *v1.Service) (runtime.Object, error)
	Remove(obj *v1.Service) (runtime.Object, error)
	Updated(obj *v1.Service) (runtime.Object, error)
}

type ServiceLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v1.Service) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v1.Service) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v1.Service) (runtime.Object, error)
}

type serviceLifecycleAdapter struct {
	lifecycle ServiceLifecycleContext
}

func (w *serviceLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *serviceLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *serviceLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *serviceLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v1.Service))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *serviceLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *serviceLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v1.Service))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *serviceLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *serviceLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v1.Service))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewServiceLifecycleAdapter(name string, clusterScoped bool, client ServiceInterface, l ServiceLifecycle) ServiceHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(ServiceGroupVersionResource)
	}
	adapter := &serviceLifecycleAdapter{lifecycle: &serviceLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v1.Service) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewServiceLifecycleAdapterContext(name string, clusterScoped bool, client ServiceInterface, l ServiceLifecycleContext) ServiceHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(ServiceGroupVersionResource)
	}
	adapter := &serviceLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v1.Service) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

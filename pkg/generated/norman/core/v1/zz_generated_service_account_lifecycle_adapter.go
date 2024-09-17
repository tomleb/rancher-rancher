package v1

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type serviceAccountLifecycleConverter struct {
	lifecycle ServiceAccountLifecycle
}

func (w *serviceAccountLifecycleConverter) CreateContext(_ context.Context, obj *v1.ServiceAccount) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *serviceAccountLifecycleConverter) RemoveContext(_ context.Context, obj *v1.ServiceAccount) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *serviceAccountLifecycleConverter) UpdatedContext(_ context.Context, obj *v1.ServiceAccount) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type ServiceAccountLifecycle interface {
	Create(obj *v1.ServiceAccount) (runtime.Object, error)
	Remove(obj *v1.ServiceAccount) (runtime.Object, error)
	Updated(obj *v1.ServiceAccount) (runtime.Object, error)
}

type ServiceAccountLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v1.ServiceAccount) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v1.ServiceAccount) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v1.ServiceAccount) (runtime.Object, error)
}

type serviceAccountLifecycleAdapter struct {
	lifecycle ServiceAccountLifecycleContext
}

func (w *serviceAccountLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *serviceAccountLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *serviceAccountLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *serviceAccountLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v1.ServiceAccount))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *serviceAccountLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *serviceAccountLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v1.ServiceAccount))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *serviceAccountLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *serviceAccountLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v1.ServiceAccount))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewServiceAccountLifecycleAdapter(name string, clusterScoped bool, client ServiceAccountInterface, l ServiceAccountLifecycle) ServiceAccountHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(ServiceAccountGroupVersionResource)
	}
	adapter := &serviceAccountLifecycleAdapter{lifecycle: &serviceAccountLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v1.ServiceAccount) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewServiceAccountLifecycleAdapterContext(name string, clusterScoped bool, client ServiceAccountInterface, l ServiceAccountLifecycleContext) ServiceAccountHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(ServiceAccountGroupVersionResource)
	}
	adapter := &serviceAccountLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v1.ServiceAccount) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

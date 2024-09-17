package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/project.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type serviceAccountTokenLifecycleConverter struct {
	lifecycle ServiceAccountTokenLifecycle
}

func (w *serviceAccountTokenLifecycleConverter) CreateContext(_ context.Context, obj *v3.ServiceAccountToken) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *serviceAccountTokenLifecycleConverter) RemoveContext(_ context.Context, obj *v3.ServiceAccountToken) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *serviceAccountTokenLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.ServiceAccountToken) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type ServiceAccountTokenLifecycle interface {
	Create(obj *v3.ServiceAccountToken) (runtime.Object, error)
	Remove(obj *v3.ServiceAccountToken) (runtime.Object, error)
	Updated(obj *v3.ServiceAccountToken) (runtime.Object, error)
}

type ServiceAccountTokenLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.ServiceAccountToken) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.ServiceAccountToken) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.ServiceAccountToken) (runtime.Object, error)
}

type serviceAccountTokenLifecycleAdapter struct {
	lifecycle ServiceAccountTokenLifecycleContext
}

func (w *serviceAccountTokenLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *serviceAccountTokenLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *serviceAccountTokenLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *serviceAccountTokenLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.ServiceAccountToken))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *serviceAccountTokenLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *serviceAccountTokenLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.ServiceAccountToken))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *serviceAccountTokenLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *serviceAccountTokenLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.ServiceAccountToken))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewServiceAccountTokenLifecycleAdapter(name string, clusterScoped bool, client ServiceAccountTokenInterface, l ServiceAccountTokenLifecycle) ServiceAccountTokenHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(ServiceAccountTokenGroupVersionResource)
	}
	adapter := &serviceAccountTokenLifecycleAdapter{lifecycle: &serviceAccountTokenLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.ServiceAccountToken) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewServiceAccountTokenLifecycleAdapterContext(name string, clusterScoped bool, client ServiceAccountTokenInterface, l ServiceAccountTokenLifecycleContext) ServiceAccountTokenHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(ServiceAccountTokenGroupVersionResource)
	}
	adapter := &serviceAccountTokenLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.ServiceAccountToken) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

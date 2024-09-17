package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/project.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type basicAuthLifecycleConverter struct {
	lifecycle BasicAuthLifecycle
}

func (w *basicAuthLifecycleConverter) CreateContext(_ context.Context, obj *v3.BasicAuth) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *basicAuthLifecycleConverter) RemoveContext(_ context.Context, obj *v3.BasicAuth) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *basicAuthLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.BasicAuth) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type BasicAuthLifecycle interface {
	Create(obj *v3.BasicAuth) (runtime.Object, error)
	Remove(obj *v3.BasicAuth) (runtime.Object, error)
	Updated(obj *v3.BasicAuth) (runtime.Object, error)
}

type BasicAuthLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.BasicAuth) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.BasicAuth) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.BasicAuth) (runtime.Object, error)
}

type basicAuthLifecycleAdapter struct {
	lifecycle BasicAuthLifecycleContext
}

func (w *basicAuthLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *basicAuthLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *basicAuthLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *basicAuthLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.BasicAuth))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *basicAuthLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *basicAuthLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.BasicAuth))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *basicAuthLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *basicAuthLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.BasicAuth))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewBasicAuthLifecycleAdapter(name string, clusterScoped bool, client BasicAuthInterface, l BasicAuthLifecycle) BasicAuthHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(BasicAuthGroupVersionResource)
	}
	adapter := &basicAuthLifecycleAdapter{lifecycle: &basicAuthLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.BasicAuth) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewBasicAuthLifecycleAdapterContext(name string, clusterScoped bool, client BasicAuthInterface, l BasicAuthLifecycleContext) BasicAuthHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(BasicAuthGroupVersionResource)
	}
	adapter := &basicAuthLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.BasicAuth) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

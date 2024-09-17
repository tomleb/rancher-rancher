package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type rkeK8sServiceOptionLifecycleConverter struct {
	lifecycle RkeK8sServiceOptionLifecycle
}

func (w *rkeK8sServiceOptionLifecycleConverter) CreateContext(_ context.Context, obj *v3.RkeK8sServiceOption) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *rkeK8sServiceOptionLifecycleConverter) RemoveContext(_ context.Context, obj *v3.RkeK8sServiceOption) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *rkeK8sServiceOptionLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.RkeK8sServiceOption) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type RkeK8sServiceOptionLifecycle interface {
	Create(obj *v3.RkeK8sServiceOption) (runtime.Object, error)
	Remove(obj *v3.RkeK8sServiceOption) (runtime.Object, error)
	Updated(obj *v3.RkeK8sServiceOption) (runtime.Object, error)
}

type RkeK8sServiceOptionLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.RkeK8sServiceOption) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.RkeK8sServiceOption) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.RkeK8sServiceOption) (runtime.Object, error)
}

type rkeK8sServiceOptionLifecycleAdapter struct {
	lifecycle RkeK8sServiceOptionLifecycleContext
}

func (w *rkeK8sServiceOptionLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *rkeK8sServiceOptionLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *rkeK8sServiceOptionLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *rkeK8sServiceOptionLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.RkeK8sServiceOption))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *rkeK8sServiceOptionLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *rkeK8sServiceOptionLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.RkeK8sServiceOption))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *rkeK8sServiceOptionLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *rkeK8sServiceOptionLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.RkeK8sServiceOption))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewRkeK8sServiceOptionLifecycleAdapter(name string, clusterScoped bool, client RkeK8sServiceOptionInterface, l RkeK8sServiceOptionLifecycle) RkeK8sServiceOptionHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(RkeK8sServiceOptionGroupVersionResource)
	}
	adapter := &rkeK8sServiceOptionLifecycleAdapter{lifecycle: &rkeK8sServiceOptionLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.RkeK8sServiceOption) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewRkeK8sServiceOptionLifecycleAdapterContext(name string, clusterScoped bool, client RkeK8sServiceOptionInterface, l RkeK8sServiceOptionLifecycleContext) RkeK8sServiceOptionHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(RkeK8sServiceOptionGroupVersionResource)
	}
	adapter := &rkeK8sServiceOptionLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.RkeK8sServiceOption) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type rkeK8sSystemImageLifecycleConverter struct {
	lifecycle RkeK8sSystemImageLifecycle
}

func (w *rkeK8sSystemImageLifecycleConverter) CreateContext(_ context.Context, obj *v3.RkeK8sSystemImage) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *rkeK8sSystemImageLifecycleConverter) RemoveContext(_ context.Context, obj *v3.RkeK8sSystemImage) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *rkeK8sSystemImageLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.RkeK8sSystemImage) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type RkeK8sSystemImageLifecycle interface {
	Create(obj *v3.RkeK8sSystemImage) (runtime.Object, error)
	Remove(obj *v3.RkeK8sSystemImage) (runtime.Object, error)
	Updated(obj *v3.RkeK8sSystemImage) (runtime.Object, error)
}

type RkeK8sSystemImageLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.RkeK8sSystemImage) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.RkeK8sSystemImage) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.RkeK8sSystemImage) (runtime.Object, error)
}

type rkeK8sSystemImageLifecycleAdapter struct {
	lifecycle RkeK8sSystemImageLifecycleContext
}

func (w *rkeK8sSystemImageLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *rkeK8sSystemImageLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *rkeK8sSystemImageLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *rkeK8sSystemImageLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.RkeK8sSystemImage))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *rkeK8sSystemImageLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *rkeK8sSystemImageLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.RkeK8sSystemImage))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *rkeK8sSystemImageLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *rkeK8sSystemImageLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.RkeK8sSystemImage))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewRkeK8sSystemImageLifecycleAdapter(name string, clusterScoped bool, client RkeK8sSystemImageInterface, l RkeK8sSystemImageLifecycle) RkeK8sSystemImageHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(RkeK8sSystemImageGroupVersionResource)
	}
	adapter := &rkeK8sSystemImageLifecycleAdapter{lifecycle: &rkeK8sSystemImageLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.RkeK8sSystemImage) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewRkeK8sSystemImageLifecycleAdapterContext(name string, clusterScoped bool, client RkeK8sSystemImageInterface, l RkeK8sSystemImageLifecycleContext) RkeK8sSystemImageHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(RkeK8sSystemImageGroupVersionResource)
	}
	adapter := &rkeK8sSystemImageLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.RkeK8sSystemImage) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

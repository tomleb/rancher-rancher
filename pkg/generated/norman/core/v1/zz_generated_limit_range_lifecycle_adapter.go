package v1

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type limitRangeLifecycleConverter struct {
	lifecycle LimitRangeLifecycle
}

func (w *limitRangeLifecycleConverter) CreateContext(_ context.Context, obj *v1.LimitRange) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *limitRangeLifecycleConverter) RemoveContext(_ context.Context, obj *v1.LimitRange) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *limitRangeLifecycleConverter) UpdatedContext(_ context.Context, obj *v1.LimitRange) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type LimitRangeLifecycle interface {
	Create(obj *v1.LimitRange) (runtime.Object, error)
	Remove(obj *v1.LimitRange) (runtime.Object, error)
	Updated(obj *v1.LimitRange) (runtime.Object, error)
}

type LimitRangeLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v1.LimitRange) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v1.LimitRange) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v1.LimitRange) (runtime.Object, error)
}

type limitRangeLifecycleAdapter struct {
	lifecycle LimitRangeLifecycleContext
}

func (w *limitRangeLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *limitRangeLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *limitRangeLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *limitRangeLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v1.LimitRange))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *limitRangeLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *limitRangeLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v1.LimitRange))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *limitRangeLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *limitRangeLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v1.LimitRange))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewLimitRangeLifecycleAdapter(name string, clusterScoped bool, client LimitRangeInterface, l LimitRangeLifecycle) LimitRangeHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(LimitRangeGroupVersionResource)
	}
	adapter := &limitRangeLifecycleAdapter{lifecycle: &limitRangeLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v1.LimitRange) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewLimitRangeLifecycleAdapterContext(name string, clusterScoped bool, client LimitRangeInterface, l LimitRangeLifecycleContext) LimitRangeHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(LimitRangeGroupVersionResource)
	}
	adapter := &limitRangeLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v1.LimitRange) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

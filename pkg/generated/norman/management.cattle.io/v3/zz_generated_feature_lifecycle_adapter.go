package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type featureLifecycleConverter struct {
	lifecycle FeatureLifecycle
}

func (w *featureLifecycleConverter) CreateContext(_ context.Context, obj *v3.Feature) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *featureLifecycleConverter) RemoveContext(_ context.Context, obj *v3.Feature) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *featureLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.Feature) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type FeatureLifecycle interface {
	Create(obj *v3.Feature) (runtime.Object, error)
	Remove(obj *v3.Feature) (runtime.Object, error)
	Updated(obj *v3.Feature) (runtime.Object, error)
}

type FeatureLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.Feature) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.Feature) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.Feature) (runtime.Object, error)
}

type featureLifecycleAdapter struct {
	lifecycle FeatureLifecycleContext
}

func (w *featureLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *featureLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *featureLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *featureLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.Feature))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *featureLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *featureLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.Feature))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *featureLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *featureLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.Feature))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewFeatureLifecycleAdapter(name string, clusterScoped bool, client FeatureInterface, l FeatureLifecycle) FeatureHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(FeatureGroupVersionResource)
	}
	adapter := &featureLifecycleAdapter{lifecycle: &featureLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.Feature) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewFeatureLifecycleAdapterContext(name string, clusterScoped bool, client FeatureInterface, l FeatureLifecycleContext) FeatureHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(FeatureGroupVersionResource)
	}
	adapter := &featureLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.Feature) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

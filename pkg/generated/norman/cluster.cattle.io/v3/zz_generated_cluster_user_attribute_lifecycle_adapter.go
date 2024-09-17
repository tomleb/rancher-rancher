package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/cluster.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type clusterUserAttributeLifecycleConverter struct {
	lifecycle ClusterUserAttributeLifecycle
}

func (w *clusterUserAttributeLifecycleConverter) CreateContext(_ context.Context, obj *v3.ClusterUserAttribute) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *clusterUserAttributeLifecycleConverter) RemoveContext(_ context.Context, obj *v3.ClusterUserAttribute) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *clusterUserAttributeLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.ClusterUserAttribute) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type ClusterUserAttributeLifecycle interface {
	Create(obj *v3.ClusterUserAttribute) (runtime.Object, error)
	Remove(obj *v3.ClusterUserAttribute) (runtime.Object, error)
	Updated(obj *v3.ClusterUserAttribute) (runtime.Object, error)
}

type ClusterUserAttributeLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.ClusterUserAttribute) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.ClusterUserAttribute) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.ClusterUserAttribute) (runtime.Object, error)
}

type clusterUserAttributeLifecycleAdapter struct {
	lifecycle ClusterUserAttributeLifecycleContext
}

func (w *clusterUserAttributeLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *clusterUserAttributeLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *clusterUserAttributeLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *clusterUserAttributeLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.ClusterUserAttribute))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *clusterUserAttributeLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *clusterUserAttributeLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.ClusterUserAttribute))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *clusterUserAttributeLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *clusterUserAttributeLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.ClusterUserAttribute))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewClusterUserAttributeLifecycleAdapter(name string, clusterScoped bool, client ClusterUserAttributeInterface, l ClusterUserAttributeLifecycle) ClusterUserAttributeHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(ClusterUserAttributeGroupVersionResource)
	}
	adapter := &clusterUserAttributeLifecycleAdapter{lifecycle: &clusterUserAttributeLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.ClusterUserAttribute) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewClusterUserAttributeLifecycleAdapterContext(name string, clusterScoped bool, client ClusterUserAttributeInterface, l ClusterUserAttributeLifecycleContext) ClusterUserAttributeHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(ClusterUserAttributeGroupVersionResource)
	}
	adapter := &clusterUserAttributeLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.ClusterUserAttribute) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

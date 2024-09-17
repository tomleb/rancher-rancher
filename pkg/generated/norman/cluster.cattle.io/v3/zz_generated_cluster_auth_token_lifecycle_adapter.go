package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/cluster.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type clusterAuthTokenLifecycleConverter struct {
	lifecycle ClusterAuthTokenLifecycle
}

func (w *clusterAuthTokenLifecycleConverter) CreateContext(_ context.Context, obj *v3.ClusterAuthToken) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *clusterAuthTokenLifecycleConverter) RemoveContext(_ context.Context, obj *v3.ClusterAuthToken) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *clusterAuthTokenLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.ClusterAuthToken) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type ClusterAuthTokenLifecycle interface {
	Create(obj *v3.ClusterAuthToken) (runtime.Object, error)
	Remove(obj *v3.ClusterAuthToken) (runtime.Object, error)
	Updated(obj *v3.ClusterAuthToken) (runtime.Object, error)
}

type ClusterAuthTokenLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.ClusterAuthToken) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.ClusterAuthToken) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.ClusterAuthToken) (runtime.Object, error)
}

type clusterAuthTokenLifecycleAdapter struct {
	lifecycle ClusterAuthTokenLifecycleContext
}

func (w *clusterAuthTokenLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *clusterAuthTokenLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *clusterAuthTokenLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *clusterAuthTokenLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.ClusterAuthToken))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *clusterAuthTokenLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *clusterAuthTokenLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.ClusterAuthToken))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *clusterAuthTokenLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *clusterAuthTokenLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.ClusterAuthToken))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewClusterAuthTokenLifecycleAdapter(name string, clusterScoped bool, client ClusterAuthTokenInterface, l ClusterAuthTokenLifecycle) ClusterAuthTokenHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(ClusterAuthTokenGroupVersionResource)
	}
	adapter := &clusterAuthTokenLifecycleAdapter{lifecycle: &clusterAuthTokenLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.ClusterAuthToken) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewClusterAuthTokenLifecycleAdapterContext(name string, clusterScoped bool, client ClusterAuthTokenInterface, l ClusterAuthTokenLifecycleContext) ClusterAuthTokenHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(ClusterAuthTokenGroupVersionResource)
	}
	adapter := &clusterAuthTokenLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.ClusterAuthToken) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

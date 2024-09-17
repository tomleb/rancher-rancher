package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type clusterRegistrationTokenLifecycleConverter struct {
	lifecycle ClusterRegistrationTokenLifecycle
}

func (w *clusterRegistrationTokenLifecycleConverter) CreateContext(_ context.Context, obj *v3.ClusterRegistrationToken) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *clusterRegistrationTokenLifecycleConverter) RemoveContext(_ context.Context, obj *v3.ClusterRegistrationToken) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *clusterRegistrationTokenLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.ClusterRegistrationToken) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type ClusterRegistrationTokenLifecycle interface {
	Create(obj *v3.ClusterRegistrationToken) (runtime.Object, error)
	Remove(obj *v3.ClusterRegistrationToken) (runtime.Object, error)
	Updated(obj *v3.ClusterRegistrationToken) (runtime.Object, error)
}

type ClusterRegistrationTokenLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.ClusterRegistrationToken) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.ClusterRegistrationToken) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.ClusterRegistrationToken) (runtime.Object, error)
}

type clusterRegistrationTokenLifecycleAdapter struct {
	lifecycle ClusterRegistrationTokenLifecycleContext
}

func (w *clusterRegistrationTokenLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *clusterRegistrationTokenLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *clusterRegistrationTokenLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *clusterRegistrationTokenLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.ClusterRegistrationToken))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *clusterRegistrationTokenLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *clusterRegistrationTokenLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.ClusterRegistrationToken))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *clusterRegistrationTokenLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *clusterRegistrationTokenLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.ClusterRegistrationToken))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewClusterRegistrationTokenLifecycleAdapter(name string, clusterScoped bool, client ClusterRegistrationTokenInterface, l ClusterRegistrationTokenLifecycle) ClusterRegistrationTokenHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(ClusterRegistrationTokenGroupVersionResource)
	}
	adapter := &clusterRegistrationTokenLifecycleAdapter{lifecycle: &clusterRegistrationTokenLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.ClusterRegistrationToken) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewClusterRegistrationTokenLifecycleAdapterContext(name string, clusterScoped bool, client ClusterRegistrationTokenInterface, l ClusterRegistrationTokenLifecycleContext) ClusterRegistrationTokenHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(ClusterRegistrationTokenGroupVersionResource)
	}
	adapter := &clusterRegistrationTokenLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.ClusterRegistrationToken) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

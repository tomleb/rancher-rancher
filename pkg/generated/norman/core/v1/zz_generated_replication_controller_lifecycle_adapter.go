package v1

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type replicationControllerLifecycleConverter struct {
	lifecycle ReplicationControllerLifecycle
}

func (w *replicationControllerLifecycleConverter) CreateContext(_ context.Context, obj *v1.ReplicationController) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *replicationControllerLifecycleConverter) RemoveContext(_ context.Context, obj *v1.ReplicationController) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *replicationControllerLifecycleConverter) UpdatedContext(_ context.Context, obj *v1.ReplicationController) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type ReplicationControllerLifecycle interface {
	Create(obj *v1.ReplicationController) (runtime.Object, error)
	Remove(obj *v1.ReplicationController) (runtime.Object, error)
	Updated(obj *v1.ReplicationController) (runtime.Object, error)
}

type ReplicationControllerLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v1.ReplicationController) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v1.ReplicationController) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v1.ReplicationController) (runtime.Object, error)
}

type replicationControllerLifecycleAdapter struct {
	lifecycle ReplicationControllerLifecycleContext
}

func (w *replicationControllerLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *replicationControllerLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *replicationControllerLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *replicationControllerLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v1.ReplicationController))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *replicationControllerLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *replicationControllerLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v1.ReplicationController))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *replicationControllerLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *replicationControllerLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v1.ReplicationController))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewReplicationControllerLifecycleAdapter(name string, clusterScoped bool, client ReplicationControllerInterface, l ReplicationControllerLifecycle) ReplicationControllerHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(ReplicationControllerGroupVersionResource)
	}
	adapter := &replicationControllerLifecycleAdapter{lifecycle: &replicationControllerLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v1.ReplicationController) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewReplicationControllerLifecycleAdapterContext(name string, clusterScoped bool, client ReplicationControllerInterface, l ReplicationControllerLifecycleContext) ReplicationControllerHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(ReplicationControllerGroupVersionResource)
	}
	adapter := &replicationControllerLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v1.ReplicationController) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

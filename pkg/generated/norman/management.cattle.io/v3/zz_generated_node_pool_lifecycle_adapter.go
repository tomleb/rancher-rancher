package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type nodePoolLifecycleConverter struct {
	lifecycle NodePoolLifecycle
}

func (w *nodePoolLifecycleConverter) CreateContext(_ context.Context, obj *v3.NodePool) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *nodePoolLifecycleConverter) RemoveContext(_ context.Context, obj *v3.NodePool) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *nodePoolLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.NodePool) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type NodePoolLifecycle interface {
	Create(obj *v3.NodePool) (runtime.Object, error)
	Remove(obj *v3.NodePool) (runtime.Object, error)
	Updated(obj *v3.NodePool) (runtime.Object, error)
}

type NodePoolLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.NodePool) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.NodePool) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.NodePool) (runtime.Object, error)
}

type nodePoolLifecycleAdapter struct {
	lifecycle NodePoolLifecycleContext
}

func (w *nodePoolLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *nodePoolLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *nodePoolLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *nodePoolLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.NodePool))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *nodePoolLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *nodePoolLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.NodePool))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *nodePoolLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *nodePoolLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.NodePool))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewNodePoolLifecycleAdapter(name string, clusterScoped bool, client NodePoolInterface, l NodePoolLifecycle) NodePoolHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(NodePoolGroupVersionResource)
	}
	adapter := &nodePoolLifecycleAdapter{lifecycle: &nodePoolLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.NodePool) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewNodePoolLifecycleAdapterContext(name string, clusterScoped bool, client NodePoolInterface, l NodePoolLifecycleContext) NodePoolHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(NodePoolGroupVersionResource)
	}
	adapter := &nodePoolLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.NodePool) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

package v1

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type networkPolicyLifecycleConverter struct {
	lifecycle NetworkPolicyLifecycle
}

func (w *networkPolicyLifecycleConverter) CreateContext(_ context.Context, obj *v1.NetworkPolicy) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *networkPolicyLifecycleConverter) RemoveContext(_ context.Context, obj *v1.NetworkPolicy) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *networkPolicyLifecycleConverter) UpdatedContext(_ context.Context, obj *v1.NetworkPolicy) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type NetworkPolicyLifecycle interface {
	Create(obj *v1.NetworkPolicy) (runtime.Object, error)
	Remove(obj *v1.NetworkPolicy) (runtime.Object, error)
	Updated(obj *v1.NetworkPolicy) (runtime.Object, error)
}

type NetworkPolicyLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v1.NetworkPolicy) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v1.NetworkPolicy) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v1.NetworkPolicy) (runtime.Object, error)
}

type networkPolicyLifecycleAdapter struct {
	lifecycle NetworkPolicyLifecycleContext
}

func (w *networkPolicyLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *networkPolicyLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *networkPolicyLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *networkPolicyLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v1.NetworkPolicy))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *networkPolicyLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *networkPolicyLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v1.NetworkPolicy))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *networkPolicyLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *networkPolicyLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v1.NetworkPolicy))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewNetworkPolicyLifecycleAdapter(name string, clusterScoped bool, client NetworkPolicyInterface, l NetworkPolicyLifecycle) NetworkPolicyHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(NetworkPolicyGroupVersionResource)
	}
	adapter := &networkPolicyLifecycleAdapter{lifecycle: &networkPolicyLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v1.NetworkPolicy) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewNetworkPolicyLifecycleAdapterContext(name string, clusterScoped bool, client NetworkPolicyInterface, l NetworkPolicyLifecycleContext) NetworkPolicyHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(NetworkPolicyGroupVersionResource)
	}
	adapter := &networkPolicyLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v1.NetworkPolicy) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

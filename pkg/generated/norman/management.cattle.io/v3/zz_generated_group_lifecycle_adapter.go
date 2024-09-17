package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type groupLifecycleConverter struct {
	lifecycle GroupLifecycle
}

func (w *groupLifecycleConverter) CreateContext(_ context.Context, obj *v3.Group) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *groupLifecycleConverter) RemoveContext(_ context.Context, obj *v3.Group) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *groupLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.Group) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type GroupLifecycle interface {
	Create(obj *v3.Group) (runtime.Object, error)
	Remove(obj *v3.Group) (runtime.Object, error)
	Updated(obj *v3.Group) (runtime.Object, error)
}

type GroupLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.Group) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.Group) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.Group) (runtime.Object, error)
}

type groupLifecycleAdapter struct {
	lifecycle GroupLifecycleContext
}

func (w *groupLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *groupLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *groupLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *groupLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.Group))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *groupLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *groupLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.Group))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *groupLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *groupLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.Group))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewGroupLifecycleAdapter(name string, clusterScoped bool, client GroupInterface, l GroupLifecycle) GroupHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(GroupGroupVersionResource)
	}
	adapter := &groupLifecycleAdapter{lifecycle: &groupLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.Group) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewGroupLifecycleAdapterContext(name string, clusterScoped bool, client GroupInterface, l GroupLifecycleContext) GroupHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(GroupGroupVersionResource)
	}
	adapter := &groupLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.Group) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

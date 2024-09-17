package v1

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type eventLifecycleConverter struct {
	lifecycle EventLifecycle
}

func (w *eventLifecycleConverter) CreateContext(_ context.Context, obj *v1.Event) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *eventLifecycleConverter) RemoveContext(_ context.Context, obj *v1.Event) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *eventLifecycleConverter) UpdatedContext(_ context.Context, obj *v1.Event) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type EventLifecycle interface {
	Create(obj *v1.Event) (runtime.Object, error)
	Remove(obj *v1.Event) (runtime.Object, error)
	Updated(obj *v1.Event) (runtime.Object, error)
}

type EventLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v1.Event) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v1.Event) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v1.Event) (runtime.Object, error)
}

type eventLifecycleAdapter struct {
	lifecycle EventLifecycleContext
}

func (w *eventLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *eventLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *eventLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *eventLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v1.Event))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *eventLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *eventLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v1.Event))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *eventLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *eventLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v1.Event))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewEventLifecycleAdapter(name string, clusterScoped bool, client EventInterface, l EventLifecycle) EventHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(EventGroupVersionResource)
	}
	adapter := &eventLifecycleAdapter{lifecycle: &eventLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v1.Event) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewEventLifecycleAdapterContext(name string, clusterScoped bool, client EventInterface, l EventLifecycleContext) EventHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(EventGroupVersionResource)
	}
	adapter := &eventLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v1.Event) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

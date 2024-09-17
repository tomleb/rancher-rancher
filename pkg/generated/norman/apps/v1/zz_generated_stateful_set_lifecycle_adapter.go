package v1

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type statefulSetLifecycleConverter struct {
	lifecycle StatefulSetLifecycle
}

func (w *statefulSetLifecycleConverter) CreateContext(_ context.Context, obj *v1.StatefulSet) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *statefulSetLifecycleConverter) RemoveContext(_ context.Context, obj *v1.StatefulSet) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *statefulSetLifecycleConverter) UpdatedContext(_ context.Context, obj *v1.StatefulSet) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type StatefulSetLifecycle interface {
	Create(obj *v1.StatefulSet) (runtime.Object, error)
	Remove(obj *v1.StatefulSet) (runtime.Object, error)
	Updated(obj *v1.StatefulSet) (runtime.Object, error)
}

type StatefulSetLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v1.StatefulSet) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v1.StatefulSet) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v1.StatefulSet) (runtime.Object, error)
}

type statefulSetLifecycleAdapter struct {
	lifecycle StatefulSetLifecycleContext
}

func (w *statefulSetLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *statefulSetLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *statefulSetLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *statefulSetLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v1.StatefulSet))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *statefulSetLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *statefulSetLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v1.StatefulSet))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *statefulSetLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *statefulSetLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v1.StatefulSet))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewStatefulSetLifecycleAdapter(name string, clusterScoped bool, client StatefulSetInterface, l StatefulSetLifecycle) StatefulSetHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(StatefulSetGroupVersionResource)
	}
	adapter := &statefulSetLifecycleAdapter{lifecycle: &statefulSetLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v1.StatefulSet) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewStatefulSetLifecycleAdapterContext(name string, clusterScoped bool, client StatefulSetInterface, l StatefulSetLifecycleContext) StatefulSetHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(StatefulSetGroupVersionResource)
	}
	adapter := &statefulSetLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v1.StatefulSet) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

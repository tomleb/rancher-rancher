package v1

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type jobLifecycleConverter struct {
	lifecycle JobLifecycle
}

func (w *jobLifecycleConverter) CreateContext(_ context.Context, obj *v1.Job) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *jobLifecycleConverter) RemoveContext(_ context.Context, obj *v1.Job) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *jobLifecycleConverter) UpdatedContext(_ context.Context, obj *v1.Job) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type JobLifecycle interface {
	Create(obj *v1.Job) (runtime.Object, error)
	Remove(obj *v1.Job) (runtime.Object, error)
	Updated(obj *v1.Job) (runtime.Object, error)
}

type JobLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v1.Job) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v1.Job) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v1.Job) (runtime.Object, error)
}

type jobLifecycleAdapter struct {
	lifecycle JobLifecycleContext
}

func (w *jobLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *jobLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *jobLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *jobLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v1.Job))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *jobLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *jobLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v1.Job))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *jobLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *jobLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v1.Job))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewJobLifecycleAdapter(name string, clusterScoped bool, client JobInterface, l JobLifecycle) JobHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(JobGroupVersionResource)
	}
	adapter := &jobLifecycleAdapter{lifecycle: &jobLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v1.Job) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewJobLifecycleAdapterContext(name string, clusterScoped bool, client JobInterface, l JobLifecycleContext) JobHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(JobGroupVersionResource)
	}
	adapter := &jobLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v1.Job) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

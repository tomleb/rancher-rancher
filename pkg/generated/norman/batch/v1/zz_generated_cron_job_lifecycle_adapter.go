package v1

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type cronJobLifecycleConverter struct {
	lifecycle CronJobLifecycle
}

func (w *cronJobLifecycleConverter) CreateContext(_ context.Context, obj *v1.CronJob) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *cronJobLifecycleConverter) RemoveContext(_ context.Context, obj *v1.CronJob) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *cronJobLifecycleConverter) UpdatedContext(_ context.Context, obj *v1.CronJob) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type CronJobLifecycle interface {
	Create(obj *v1.CronJob) (runtime.Object, error)
	Remove(obj *v1.CronJob) (runtime.Object, error)
	Updated(obj *v1.CronJob) (runtime.Object, error)
}

type CronJobLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v1.CronJob) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v1.CronJob) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v1.CronJob) (runtime.Object, error)
}

type cronJobLifecycleAdapter struct {
	lifecycle CronJobLifecycleContext
}

func (w *cronJobLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *cronJobLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *cronJobLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *cronJobLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v1.CronJob))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *cronJobLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *cronJobLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v1.CronJob))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *cronJobLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *cronJobLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v1.CronJob))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewCronJobLifecycleAdapter(name string, clusterScoped bool, client CronJobInterface, l CronJobLifecycle) CronJobHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(CronJobGroupVersionResource)
	}
	adapter := &cronJobLifecycleAdapter{lifecycle: &cronJobLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v1.CronJob) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewCronJobLifecycleAdapterContext(name string, clusterScoped bool, client CronJobInterface, l CronJobLifecycleContext) CronJobHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(CronJobGroupVersionResource)
	}
	adapter := &cronJobLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v1.CronJob) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

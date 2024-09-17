package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type projectLifecycleConverter struct {
	lifecycle ProjectLifecycle
}

func (w *projectLifecycleConverter) CreateContext(_ context.Context, obj *v3.Project) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *projectLifecycleConverter) RemoveContext(_ context.Context, obj *v3.Project) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *projectLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.Project) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type ProjectLifecycle interface {
	Create(obj *v3.Project) (runtime.Object, error)
	Remove(obj *v3.Project) (runtime.Object, error)
	Updated(obj *v3.Project) (runtime.Object, error)
}

type ProjectLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.Project) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.Project) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.Project) (runtime.Object, error)
}

type projectLifecycleAdapter struct {
	lifecycle ProjectLifecycleContext
}

func (w *projectLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *projectLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *projectLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *projectLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.Project))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *projectLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *projectLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.Project))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *projectLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *projectLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.Project))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewProjectLifecycleAdapter(name string, clusterScoped bool, client ProjectInterface, l ProjectLifecycle) ProjectHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(ProjectGroupVersionResource)
	}
	adapter := &projectLifecycleAdapter{lifecycle: &projectLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.Project) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewProjectLifecycleAdapterContext(name string, clusterScoped bool, client ProjectInterface, l ProjectLifecycleContext) ProjectHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(ProjectGroupVersionResource)
	}
	adapter := &projectLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.Project) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type principalLifecycleConverter struct {
	lifecycle PrincipalLifecycle
}

func (w *principalLifecycleConverter) CreateContext(_ context.Context, obj *v3.Principal) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *principalLifecycleConverter) RemoveContext(_ context.Context, obj *v3.Principal) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *principalLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.Principal) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type PrincipalLifecycle interface {
	Create(obj *v3.Principal) (runtime.Object, error)
	Remove(obj *v3.Principal) (runtime.Object, error)
	Updated(obj *v3.Principal) (runtime.Object, error)
}

type PrincipalLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.Principal) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.Principal) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.Principal) (runtime.Object, error)
}

type principalLifecycleAdapter struct {
	lifecycle PrincipalLifecycleContext
}

func (w *principalLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *principalLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *principalLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *principalLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.Principal))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *principalLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *principalLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.Principal))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *principalLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *principalLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.Principal))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewPrincipalLifecycleAdapter(name string, clusterScoped bool, client PrincipalInterface, l PrincipalLifecycle) PrincipalHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(PrincipalGroupVersionResource)
	}
	adapter := &principalLifecycleAdapter{lifecycle: &principalLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.Principal) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewPrincipalLifecycleAdapterContext(name string, clusterScoped bool, client PrincipalInterface, l PrincipalLifecycleContext) PrincipalHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(PrincipalGroupVersionResource)
	}
	adapter := &principalLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.Principal) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

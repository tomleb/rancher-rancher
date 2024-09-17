package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type composeConfigLifecycleConverter struct {
	lifecycle ComposeConfigLifecycle
}

func (w *composeConfigLifecycleConverter) CreateContext(_ context.Context, obj *v3.ComposeConfig) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *composeConfigLifecycleConverter) RemoveContext(_ context.Context, obj *v3.ComposeConfig) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *composeConfigLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.ComposeConfig) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type ComposeConfigLifecycle interface {
	Create(obj *v3.ComposeConfig) (runtime.Object, error)
	Remove(obj *v3.ComposeConfig) (runtime.Object, error)
	Updated(obj *v3.ComposeConfig) (runtime.Object, error)
}

type ComposeConfigLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.ComposeConfig) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.ComposeConfig) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.ComposeConfig) (runtime.Object, error)
}

type composeConfigLifecycleAdapter struct {
	lifecycle ComposeConfigLifecycleContext
}

func (w *composeConfigLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *composeConfigLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *composeConfigLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *composeConfigLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.ComposeConfig))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *composeConfigLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *composeConfigLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.ComposeConfig))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *composeConfigLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *composeConfigLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.ComposeConfig))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewComposeConfigLifecycleAdapter(name string, clusterScoped bool, client ComposeConfigInterface, l ComposeConfigLifecycle) ComposeConfigHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(ComposeConfigGroupVersionResource)
	}
	adapter := &composeConfigLifecycleAdapter{lifecycle: &composeConfigLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.ComposeConfig) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewComposeConfigLifecycleAdapterContext(name string, clusterScoped bool, client ComposeConfigInterface, l ComposeConfigLifecycleContext) ComposeConfigHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(ComposeConfigGroupVersionResource)
	}
	adapter := &composeConfigLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.ComposeConfig) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

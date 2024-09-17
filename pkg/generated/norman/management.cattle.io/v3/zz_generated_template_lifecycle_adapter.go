package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type templateLifecycleConverter struct {
	lifecycle TemplateLifecycle
}

func (w *templateLifecycleConverter) CreateContext(_ context.Context, obj *v3.Template) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *templateLifecycleConverter) RemoveContext(_ context.Context, obj *v3.Template) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *templateLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.Template) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type TemplateLifecycle interface {
	Create(obj *v3.Template) (runtime.Object, error)
	Remove(obj *v3.Template) (runtime.Object, error)
	Updated(obj *v3.Template) (runtime.Object, error)
}

type TemplateLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.Template) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.Template) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.Template) (runtime.Object, error)
}

type templateLifecycleAdapter struct {
	lifecycle TemplateLifecycleContext
}

func (w *templateLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *templateLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *templateLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *templateLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.Template))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *templateLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *templateLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.Template))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *templateLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *templateLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.Template))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewTemplateLifecycleAdapter(name string, clusterScoped bool, client TemplateInterface, l TemplateLifecycle) TemplateHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(TemplateGroupVersionResource)
	}
	adapter := &templateLifecycleAdapter{lifecycle: &templateLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.Template) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewTemplateLifecycleAdapterContext(name string, clusterScoped bool, client TemplateInterface, l TemplateLifecycleContext) TemplateHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(TemplateGroupVersionResource)
	}
	adapter := &templateLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.Template) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

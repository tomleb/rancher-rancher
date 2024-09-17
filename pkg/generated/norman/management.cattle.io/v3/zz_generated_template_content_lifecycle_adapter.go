package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type templateContentLifecycleConverter struct {
	lifecycle TemplateContentLifecycle
}

func (w *templateContentLifecycleConverter) CreateContext(_ context.Context, obj *v3.TemplateContent) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *templateContentLifecycleConverter) RemoveContext(_ context.Context, obj *v3.TemplateContent) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *templateContentLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.TemplateContent) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type TemplateContentLifecycle interface {
	Create(obj *v3.TemplateContent) (runtime.Object, error)
	Remove(obj *v3.TemplateContent) (runtime.Object, error)
	Updated(obj *v3.TemplateContent) (runtime.Object, error)
}

type TemplateContentLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.TemplateContent) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.TemplateContent) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.TemplateContent) (runtime.Object, error)
}

type templateContentLifecycleAdapter struct {
	lifecycle TemplateContentLifecycleContext
}

func (w *templateContentLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *templateContentLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *templateContentLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *templateContentLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.TemplateContent))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *templateContentLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *templateContentLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.TemplateContent))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *templateContentLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *templateContentLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.TemplateContent))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewTemplateContentLifecycleAdapter(name string, clusterScoped bool, client TemplateContentInterface, l TemplateContentLifecycle) TemplateContentHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(TemplateContentGroupVersionResource)
	}
	adapter := &templateContentLifecycleAdapter{lifecycle: &templateContentLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.TemplateContent) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewTemplateContentLifecycleAdapterContext(name string, clusterScoped bool, client TemplateContentInterface, l TemplateContentLifecycleContext) TemplateContentHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(TemplateContentGroupVersionResource)
	}
	adapter := &templateContentLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.TemplateContent) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type templateVersionLifecycleConverter struct {
	lifecycle TemplateVersionLifecycle
}

func (w *templateVersionLifecycleConverter) CreateContext(_ context.Context, obj *v3.TemplateVersion) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *templateVersionLifecycleConverter) RemoveContext(_ context.Context, obj *v3.TemplateVersion) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *templateVersionLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.TemplateVersion) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type TemplateVersionLifecycle interface {
	Create(obj *v3.TemplateVersion) (runtime.Object, error)
	Remove(obj *v3.TemplateVersion) (runtime.Object, error)
	Updated(obj *v3.TemplateVersion) (runtime.Object, error)
}

type TemplateVersionLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.TemplateVersion) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.TemplateVersion) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.TemplateVersion) (runtime.Object, error)
}

type templateVersionLifecycleAdapter struct {
	lifecycle TemplateVersionLifecycleContext
}

func (w *templateVersionLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *templateVersionLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *templateVersionLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *templateVersionLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.TemplateVersion))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *templateVersionLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *templateVersionLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.TemplateVersion))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *templateVersionLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *templateVersionLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.TemplateVersion))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewTemplateVersionLifecycleAdapter(name string, clusterScoped bool, client TemplateVersionInterface, l TemplateVersionLifecycle) TemplateVersionHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(TemplateVersionGroupVersionResource)
	}
	adapter := &templateVersionLifecycleAdapter{lifecycle: &templateVersionLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.TemplateVersion) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewTemplateVersionLifecycleAdapterContext(name string, clusterScoped bool, client TemplateVersionInterface, l TemplateVersionLifecycleContext) TemplateVersionHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(TemplateVersionGroupVersionResource)
	}
	adapter := &templateVersionLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.TemplateVersion) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

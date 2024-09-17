package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type projectRoleTemplateBindingLifecycleConverter struct {
	lifecycle ProjectRoleTemplateBindingLifecycle
}

func (w *projectRoleTemplateBindingLifecycleConverter) CreateContext(_ context.Context, obj *v3.ProjectRoleTemplateBinding) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *projectRoleTemplateBindingLifecycleConverter) RemoveContext(_ context.Context, obj *v3.ProjectRoleTemplateBinding) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *projectRoleTemplateBindingLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.ProjectRoleTemplateBinding) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type ProjectRoleTemplateBindingLifecycle interface {
	Create(obj *v3.ProjectRoleTemplateBinding) (runtime.Object, error)
	Remove(obj *v3.ProjectRoleTemplateBinding) (runtime.Object, error)
	Updated(obj *v3.ProjectRoleTemplateBinding) (runtime.Object, error)
}

type ProjectRoleTemplateBindingLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.ProjectRoleTemplateBinding) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.ProjectRoleTemplateBinding) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.ProjectRoleTemplateBinding) (runtime.Object, error)
}

type projectRoleTemplateBindingLifecycleAdapter struct {
	lifecycle ProjectRoleTemplateBindingLifecycleContext
}

func (w *projectRoleTemplateBindingLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *projectRoleTemplateBindingLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *projectRoleTemplateBindingLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *projectRoleTemplateBindingLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.ProjectRoleTemplateBinding))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *projectRoleTemplateBindingLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *projectRoleTemplateBindingLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.ProjectRoleTemplateBinding))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *projectRoleTemplateBindingLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *projectRoleTemplateBindingLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.ProjectRoleTemplateBinding))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewProjectRoleTemplateBindingLifecycleAdapter(name string, clusterScoped bool, client ProjectRoleTemplateBindingInterface, l ProjectRoleTemplateBindingLifecycle) ProjectRoleTemplateBindingHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(ProjectRoleTemplateBindingGroupVersionResource)
	}
	adapter := &projectRoleTemplateBindingLifecycleAdapter{lifecycle: &projectRoleTemplateBindingLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.ProjectRoleTemplateBinding) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewProjectRoleTemplateBindingLifecycleAdapterContext(name string, clusterScoped bool, client ProjectRoleTemplateBindingInterface, l ProjectRoleTemplateBindingLifecycleContext) ProjectRoleTemplateBindingHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(ProjectRoleTemplateBindingGroupVersionResource)
	}
	adapter := &projectRoleTemplateBindingLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.ProjectRoleTemplateBinding) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

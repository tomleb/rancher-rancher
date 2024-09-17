package v1

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type deploymentLifecycleConverter struct {
	lifecycle DeploymentLifecycle
}

func (w *deploymentLifecycleConverter) CreateContext(_ context.Context, obj *v1.Deployment) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *deploymentLifecycleConverter) RemoveContext(_ context.Context, obj *v1.Deployment) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *deploymentLifecycleConverter) UpdatedContext(_ context.Context, obj *v1.Deployment) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type DeploymentLifecycle interface {
	Create(obj *v1.Deployment) (runtime.Object, error)
	Remove(obj *v1.Deployment) (runtime.Object, error)
	Updated(obj *v1.Deployment) (runtime.Object, error)
}

type DeploymentLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v1.Deployment) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v1.Deployment) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v1.Deployment) (runtime.Object, error)
}

type deploymentLifecycleAdapter struct {
	lifecycle DeploymentLifecycleContext
}

func (w *deploymentLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *deploymentLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *deploymentLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *deploymentLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v1.Deployment))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *deploymentLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *deploymentLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v1.Deployment))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *deploymentLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *deploymentLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v1.Deployment))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewDeploymentLifecycleAdapter(name string, clusterScoped bool, client DeploymentInterface, l DeploymentLifecycle) DeploymentHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(DeploymentGroupVersionResource)
	}
	adapter := &deploymentLifecycleAdapter{lifecycle: &deploymentLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v1.Deployment) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewDeploymentLifecycleAdapterContext(name string, clusterScoped bool, client DeploymentInterface, l DeploymentLifecycleContext) DeploymentHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(DeploymentGroupVersionResource)
	}
	adapter := &deploymentLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v1.Deployment) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

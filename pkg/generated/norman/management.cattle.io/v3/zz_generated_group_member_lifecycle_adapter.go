package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type groupMemberLifecycleConverter struct {
	lifecycle GroupMemberLifecycle
}

func (w *groupMemberLifecycleConverter) CreateContext(_ context.Context, obj *v3.GroupMember) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *groupMemberLifecycleConverter) RemoveContext(_ context.Context, obj *v3.GroupMember) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *groupMemberLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.GroupMember) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type GroupMemberLifecycle interface {
	Create(obj *v3.GroupMember) (runtime.Object, error)
	Remove(obj *v3.GroupMember) (runtime.Object, error)
	Updated(obj *v3.GroupMember) (runtime.Object, error)
}

type GroupMemberLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.GroupMember) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.GroupMember) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.GroupMember) (runtime.Object, error)
}

type groupMemberLifecycleAdapter struct {
	lifecycle GroupMemberLifecycleContext
}

func (w *groupMemberLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *groupMemberLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *groupMemberLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *groupMemberLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.GroupMember))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *groupMemberLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *groupMemberLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.GroupMember))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *groupMemberLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *groupMemberLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.GroupMember))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewGroupMemberLifecycleAdapter(name string, clusterScoped bool, client GroupMemberInterface, l GroupMemberLifecycle) GroupMemberHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(GroupMemberGroupVersionResource)
	}
	adapter := &groupMemberLifecycleAdapter{lifecycle: &groupMemberLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.GroupMember) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewGroupMemberLifecycleAdapterContext(name string, clusterScoped bool, client GroupMemberInterface, l GroupMemberLifecycleContext) GroupMemberHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(GroupMemberGroupVersionResource)
	}
	adapter := &groupMemberLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.GroupMember) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

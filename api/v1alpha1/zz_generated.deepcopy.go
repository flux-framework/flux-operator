//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ContainerResource.
func (in *ContainerResource) DeepCopy() *ContainerResource {
	if in == nil {
		return nil
	}
	out := new(ContainerResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ContainerResources) DeepCopyInto(out *ContainerResources) {
	*out = *in
	in.Limits.DeepCopyInto(&out.Limits)
	in.Requests.DeepCopyInto(&out.Requests)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ContainerResources.
func (in *ContainerResources) DeepCopy() *ContainerResources {
	if in == nil {
		return nil
	}
	out := new(ContainerResources)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ContainerVolume) DeepCopyInto(out *ContainerVolume) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ContainerVolume.
func (in *ContainerVolume) DeepCopy() *ContainerVolume {
	if in == nil {
		return nil
	}
	out := new(ContainerVolume)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FluxRestful) DeepCopyInto(out *FluxRestful) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FluxRestful.
func (in *FluxRestful) DeepCopy() *FluxRestful {
	if in == nil {
		return nil
	}
	out := new(FluxRestful)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MiniCluster) DeepCopyInto(out *MiniCluster) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MiniCluster.
func (in *MiniCluster) DeepCopy() *MiniCluster {
	if in == nil {
		return nil
	}
	out := new(MiniCluster)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MiniCluster) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MiniClusterContainer) DeepCopyInto(out *MiniClusterContainer) {
	*out = *in
	if in.Ports != nil {
		in, out := &in.Ports, &out.Ports
		*out = make([]int32, len(*in))
		copy(*out, *in)
	}
	if in.Envars != nil {
		in, out := &in.Envars, &out.Envars
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Volumes != nil {
		in, out := &in.Volumes, &out.Volumes
		*out = make(map[string]ContainerVolume, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	in.Resources.DeepCopyInto(&out.Resources)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MiniClusterContainer.
func (in *MiniClusterContainer) DeepCopy() *MiniClusterContainer {
	if in == nil {
		return nil
	}
	out := new(MiniClusterContainer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MiniClusterList) DeepCopyInto(out *MiniClusterList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]MiniCluster, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MiniClusterList.
func (in *MiniClusterList) DeepCopy() *MiniClusterList {
	if in == nil {
		return nil
	}
	out := new(MiniClusterList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MiniClusterList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MiniClusterSpec) DeepCopyInto(out *MiniClusterSpec) {
	*out = *in
	if in.Containers != nil {
		in, out := &in.Containers, &out.Containers
		*out = make([]MiniClusterContainer, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.JobLabels != nil {
		in, out := &in.JobLabels, &out.JobLabels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.PodLabels != nil {
		in, out := &in.PodLabels, &out.PodLabels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Volumes != nil {
		in, out := &in.Volumes, &out.Volumes
		*out = make(map[string]MiniClusterVolume, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	out.FluxRestful = in.FluxRestful
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MiniClusterSpec.
func (in *MiniClusterSpec) DeepCopy() *MiniClusterSpec {
	if in == nil {
		return nil
	}
	out := new(MiniClusterSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MiniClusterStatus) DeepCopyInto(out *MiniClusterStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MiniClusterStatus.
func (in *MiniClusterStatus) DeepCopy() *MiniClusterStatus {
	if in == nil {
		return nil
	}
	out := new(MiniClusterStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MiniClusterVolume) DeepCopyInto(out *MiniClusterVolume) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MiniClusterVolume.
func (in *MiniClusterVolume) DeepCopy() *MiniClusterVolume {
	if in == nil {
		return nil
	}
	out := new(MiniClusterVolume)
	in.DeepCopyInto(out)
	return out
}

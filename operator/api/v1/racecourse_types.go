/*
Copyright 2023.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RacecourseSpec defines the desired state of Racecourse
type RacecourseSpec struct {
	// Important: Run "make" to regenerate code after modifying this file

	// DeploymentName is what the controller will name child resources
	DeploymentName string `json:"deploymentName,omitempty"`

	// Override image string for deployment, not actually implemented
	Image string `json:"image,omitempty"`
}

// RacecourseStatus defines the observed state of Racecourse
type RacecourseStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Racecourse is the Schema for the racecourses API
type Racecourse struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RacecourseSpec   `json:"spec,omitempty"`
	Status RacecourseStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RacecourseList contains a list of Racecourse
type RacecourseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Racecourse `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Racecourse{}, &RacecourseList{})
}

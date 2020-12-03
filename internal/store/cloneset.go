/*
Copyright 2016 The Kubernetes Authors All rights reserved.

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

package store

import (
	"k8s.io/kube-state-metrics/pkg/metric"

	clientset "github.com/SchoIsles/kruise-state-metrics/pkg/client/clientset/versioned"
	kruiseappsv1alpha1 "github.com/openkruise/kruise-api/apps/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

var (
	descCloneSetLabelsName          = "kube_cloneset_labels"
	descCloneSetLabelsHelp          = "Kubernetes labels converted to Prometheus labels."
	descCloneSetLabelsDefaultLabels = []string{"namespace", "cloneset"}

	clonesetMetricFamilies = []metric.FamilyGenerator{
		{
			Name: "kube_cloneset_created",
			Type: metric.Gauge,
			Help: "Unix creation timestamp",
			GenerateFunc: wrapCloneSetFunc(func(d *kruiseappsv1alpha1.CloneSet) *metric.Family {
				ms := []*metric.Metric{}

				if !d.CreationTimestamp.IsZero() {
					ms = append(ms, &metric.Metric{
						Value: float64(d.CreationTimestamp.Unix()),
					})
				}

				return &metric.Family{
					Metrics: ms,
				}
			}),
		},
		{
			Name: "kube_cloneset_status_replicas",
			Type: metric.Gauge,
			Help: "The number of replicas per cloneset.",
			GenerateFunc: wrapCloneSetFunc(func(d *kruiseappsv1alpha1.CloneSet) *metric.Family {
				return &metric.Family{
					Metrics: []*metric.Metric{
						{
							Value: float64(d.Status.Replicas),
						},
					},
				}
			}),
		},
		{
			Name: "kube_cloneset_status_replicas_available",
			Type: metric.Gauge,
			Help: "The number of available replicas per cloneset.",
			GenerateFunc: wrapCloneSetFunc(func(d *kruiseappsv1alpha1.CloneSet) *metric.Family {
				return &metric.Family{
					Metrics: []*metric.Metric{
						{
							Value: float64(d.Status.AvailableReplicas),
						},
					},
				}
			}),
		},
		{
			Name: "kube_cloneset_status_replicas_unavailable",
			Type: metric.Gauge,
			Help: "The number of unavailable replicas per cloneset.",
			GenerateFunc: wrapCloneSetFunc(func(d *kruiseappsv1alpha1.CloneSet) *metric.Family {
				return &metric.Family{
					Metrics: []*metric.Metric{
						{
							Value: float64(d.Status.Replicas - d.Status.AvailableReplicas),
						},
					},
				}
			}),
		},
		{
			Name: "kube_cloneset_status_replicas_updated",
			Type: metric.Gauge,
			Help: "The number of updated replicas per cloneset.",
			GenerateFunc: wrapCloneSetFunc(func(d *kruiseappsv1alpha1.CloneSet) *metric.Family {
				return &metric.Family{
					Metrics: []*metric.Metric{
						{
							Value: float64(d.Status.UpdatedReplicas),
						},
					},
				}
			}),
		},
		{
			Name: "kube_cloneset_status_replicas_ready_updated",
			Type: metric.Gauge,
			Help: "The number of ready updated replicas per cloneset.",
			GenerateFunc: wrapCloneSetFunc(func(d *kruiseappsv1alpha1.CloneSet) *metric.Family {
				return &metric.Family{
					Metrics: []*metric.Metric{
						{
							Value: float64(d.Status.UpdatedReadyReplicas),
						},
					},
				}
			}),
		},
		{
			Name: "kube_cloneset_status_observed_generation",
			Type: metric.Gauge,
			Help: "The generation observed by the cloneset controller.",
			GenerateFunc: wrapCloneSetFunc(func(d *kruiseappsv1alpha1.CloneSet) *metric.Family {
				return &metric.Family{
					Metrics: []*metric.Metric{
						{
							Value: float64(d.Status.ObservedGeneration),
						},
					},
				}
			}),
		},
		{
			Name: "kube_cloneset_status_condition",
			Type: metric.Gauge,
			Help: "The current status conditions of a cloneset.",
			GenerateFunc: wrapCloneSetFunc(func(d *kruiseappsv1alpha1.CloneSet) *metric.Family {
				ms := make([]*metric.Metric, len(d.Status.Conditions)*len(conditionStatuses))

				for i, c := range d.Status.Conditions {
					conditionMetrics := addConditionMetrics(c.Status)

					for j, m := range conditionMetrics {
						metric := m

						metric.LabelKeys = []string{"condition", "status"}
						metric.LabelValues = append([]string{string(c.Type)}, metric.LabelValues...)
						ms[i*len(conditionStatuses)+j] = metric
					}
				}

				return &metric.Family{
					Metrics: ms,
				}
			}),
		},
		{
			Name: "kube_cloneset_spec_replicas",
			Type: metric.Gauge,
			Help: "Number of desired pods for a cloneset.",
			GenerateFunc: wrapCloneSetFunc(func(d *kruiseappsv1alpha1.CloneSet) *metric.Family {
				return &metric.Family{
					Metrics: []*metric.Metric{
						{
							Value: float64(*d.Spec.Replicas),
						},
					},
				}
			}),
		},
		{
			Name: "kube_cloneset_spec_paused",
			Type: metric.Gauge,
			Help: "Whether the cloneset is paused and will not be processed by the cloneset controller.",
			GenerateFunc: wrapCloneSetFunc(func(d *kruiseappsv1alpha1.CloneSet) *metric.Family {
				return &metric.Family{
					Metrics: []*metric.Metric{
						{
							Value: boolFloat64(d.Spec.UpdateStrategy.Paused),
						},
					},
				}
			}),
		},
		{
			Name: "kube_cloneset_spec_strategy_rollingupdate_max_unavailable",
			Type: metric.Gauge,
			Help: "Maximum number of unavailable replicas during a rolling update of a cloneset.",
			GenerateFunc: wrapCloneSetFunc(func(d *kruiseappsv1alpha1.CloneSet) *metric.Family {
				if d.Spec.UpdateStrategy.MaxUnavailable == nil {
					return &metric.Family{}
				}

				maxUnavailable, err := intstr.GetValueFromIntOrPercent(d.Spec.UpdateStrategy.MaxUnavailable, int(*d.Spec.Replicas), true)
				if err != nil {
					panic(err)
				}

				return &metric.Family{
					Metrics: []*metric.Metric{
						{
							Value: float64(maxUnavailable),
						},
					},
				}
			}),
		},
		{
			Name: "kube_cloneset_spec_strategy_rollingupdate_max_surge",
			Type: metric.Gauge,
			Help: "Maximum number of replicas that can be scheduled above the desired number of replicas during a rolling update of a cloneset.",
			GenerateFunc: wrapCloneSetFunc(func(d *kruiseappsv1alpha1.CloneSet) *metric.Family {
				if d.Spec.UpdateStrategy.MaxSurge == nil {
					return &metric.Family{}
				}

				maxSurge, err := intstr.GetValueFromIntOrPercent(d.Spec.UpdateStrategy.MaxSurge, int(*d.Spec.Replicas), true)
				if err != nil {
					panic(err)
				}

				return &metric.Family{
					Metrics: []*metric.Metric{
						{
							Value: float64(maxSurge),
						},
					},
				}
			}),
		},
		{
			Name: "kube_cloneset_metadata_generation",
			Type: metric.Gauge,
			Help: "Sequence number representing a specific generation of the desired state.",
			GenerateFunc: wrapCloneSetFunc(func(d *kruiseappsv1alpha1.CloneSet) *metric.Family {
				return &metric.Family{
					Metrics: []*metric.Metric{
						{
							Value: float64(d.ObjectMeta.Generation),
						},
					},
				}
			}),
		},
		{
			Name: descCloneSetLabelsName,
			Type: metric.Gauge,
			Help: descCloneSetLabelsHelp,
			GenerateFunc: wrapCloneSetFunc(func(d *kruiseappsv1alpha1.CloneSet) *metric.Family {
				labelKeys, labelValues := kubeLabelsToPrometheusLabels(d.Labels)
				return &metric.Family{
					Metrics: []*metric.Metric{
						{
							LabelKeys:   labelKeys,
							LabelValues: labelValues,
							Value:       1,
						},
					},
				}
			}),
		},
	}
)

func wrapCloneSetFunc(f func(*kruiseappsv1alpha1.CloneSet) *metric.Family) func(interface{}) *metric.Family {
	return func(obj interface{}) *metric.Family {
		cloneset := obj.(*kruiseappsv1alpha1.CloneSet)

		metricFamily := f(cloneset)

		for _, m := range metricFamily.Metrics {
			m.LabelKeys = append(descCloneSetLabelsDefaultLabels, m.LabelKeys...)
			m.LabelValues = append([]string{cloneset.Namespace, cloneset.Name}, m.LabelValues...)
		}

		return metricFamily
	}
}

func createCloneSetListWatch(kubeClient clientset.Interface, ns string) cache.ListerWatcher {
	return &cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			return kubeClient.AppsV1alpha1().CloneSets(ns).List(opts)
		},
		WatchFunc: func(opts metav1.ListOptions) (watch.Interface, error) {
			return kubeClient.AppsV1alpha1().CloneSets(ns).Watch(opts)
		},
	}
}

/*
Copyright (C) 2022 Christoph Görn

This file is part of r-gespraech.

r-gespraech is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

r-gespraech is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with r-gespraech.  If not, see <http://www.gnu.org/licenses/>.
*/

package controllers

import (
	"fmt"
	"time"

	kbatch "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/goern/r-gespraech/api/v1alpha1"
)

var _ = Describe("CallbackUrl controller", func() {
	var (
		testNamespace string = "default"
	)
	const (
		testCallbackUrlName  = "abc123"
		testAdviserId        = "abc123"
		CallbackPayloadAName = "abc123"

		timeout  = time.Second * 10
		duration = time.Second * 10
		interval = time.Millisecond * 250
	)

	BeforeEach(func() {
		testNamespace = "test-" + String(6)
		nsSpec := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: testNamespace}}
		Expect(k8sClient.Create(ctx, nsSpec)).Should(Succeed())
	})

	Context("When creating a CallbackUrl and no CallbackPayload", func() {
		It("Should have NoAssociatedPayloads Condition", func() {

			By("By creating a new CallbackUrl")
			callbackUrl := generateCallbackUrl(testCallbackUrlName, testNamespace, "https://localhost.local:8181/webhook/xyz_callback")
			Expect(k8sClient.Create(ctx, callbackUrl)).Should(Succeed())

			By("By checking the CallbackUrl has NoAssociatedPayloads Conditions")
			lookupKey := types.NamespacedName{Name: testCallbackUrlName, Namespace: testNamespace}
			createdCallbackUrl := &v1alpha1.CallbackUrl{}

			// We'll need to retry getting this newly created CronJob, given that creation may not immediately happen.
			Eventually(func() bool {
				err := k8sClient.Get(ctx, lookupKey, createdCallbackUrl)
				return err == nil
			}, timeout, interval).Should(BeTrue())

			Expect(meta.IsStatusConditionPresentAndEqual(callbackUrl.Status.Conditions, v1alpha1.NoAssociatedPayloads, metav1.ConditionTrue))
		})
	})
	Context("When creating a CallbackUrl and one associated CallbackPayload", func() {
		var callbackUrl *v1alpha1.CallbackUrl

		It("Should have AssociatedPayloads Condition", func() {

			By("By creating a new CallbackUrl")
			callbackUrl = generateCallbackUrl(testCallbackUrlName, testNamespace, "https://localhost.local:8181/webhook/xyz_callback")
			Expect(k8sClient.Create(ctx, callbackUrl)).Should(Succeed())

			By("By checking the CallbackUrl has no associated CallbackPayloads (aka AssociatedPayloads Condition == false!)")
			lookupKey := types.NamespacedName{Name: testCallbackUrlName, Namespace: testNamespace}
			Consistently(func() (bool, error) {
				err := k8sClient.Get(ctx, lookupKey, callbackUrl)
				if err != nil {
					return false, err
				}
				return meta.IsStatusConditionPresentAndEqual(callbackUrl.Status.Conditions, v1alpha1.AssociatedPayloads, metav1.ConditionTrue), nil
			}, duration, interval).Should(BeFalse())

			By("By creating a new CallbackPayload")
			callbackPayload := generateCallbackPayload(testAdviserId, testNamespace)
			Expect(k8sClient.Create(ctx, callbackPayload)).Should(Succeed())

			By("By checking the CallbackUrl has an AssociatedPayloads Conditions set to True")
			Consistently(func() ([]metav1.Condition, error) {
				err := k8sClient.Get(ctx, lookupKey, callbackUrl)
				if err != nil {
					return nil, err
				}
				return callbackUrl.Status.Conditions, nil
			}, timeout, interval).Should(Not(BeNil()))

			Expect(meta.IsStatusConditionPresentAndEqual(callbackUrl.Status.Conditions, v1alpha1.AssociatedPayloads, metav1.ConditionTrue))

		})

		It("should also create a Job associated with the CallbackUrl**Payload combination", func() {
			By("By checking if a Job gets created for the CallbackUrl and it's CallbackPayload")
			var jobs kbatch.JobList

			err := k8sClient.List(ctx, &jobs, client.InNamespace(testNamespace), client.MatchingFields{jobOwnerKey: callbackUrl.Name})
			if err != nil {
				fmt.Printf("%v", jobs)
			}

			Expect(jobs.Items).To(Not(BeNil()))

		})
	})
})

func generateCallbackUrl(adviserId string, namespace string, url string) *v1alpha1.CallbackUrl {
	labels := make(map[string]string)
	labels["adviser.thoth-station.ninja/adviser-id"] = adviserId

	callbackUrl := &v1alpha1.CallbackUrl{
		TypeMeta:   metav1.TypeMeta{APIVersion: "erinnerung.thoth-station.ninja/v1alpha1", Kind: "CallbackUrl"},
		ObjectMeta: metav1.ObjectMeta{Name: adviserId, Namespace: namespace, Labels: labels},
		Spec: v1alpha1.CallbackUrlSpec{
			URL: url,
			Selector: metav1.LabelSelector{
				MatchLabels: labels,
			},
		},
		Status: v1alpha1.CallbackUrlStatus{},
	}
	return callbackUrl
}

func generateCallbackPayload(adviserId string, namespace string) *v1alpha1.CallbackPayload {
	labels := make(map[string]string)
	labels["adviser.thoth-station.ninja/adviser-id"] = adviserId

	callbackPayload := &v1alpha1.CallbackPayload{
		TypeMeta:   metav1.TypeMeta{APIVersion: "erinnerung.thoth-station.ninja/v1alpha1", Kind: "CallbackPayload"},
		ObjectMeta: metav1.ObjectMeta{Name: adviserId, Namespace: namespace, Labels: labels},
		Spec: v1alpha1.CallbackPayloadSpec{
			Data: "{'adviser-document-id': adviser_id}",
			Selector: metav1.LabelSelector{
				MatchLabels: labels,
			},
		},
		Status: v1alpha1.CallbackPayloadStatus{},
	}

	return callbackPayload
}

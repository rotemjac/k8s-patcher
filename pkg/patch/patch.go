package patch

import (
	"context"
	"encoding/json"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	client "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"strings"
)

type PatchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

func Patch(podName string, podNamespace string, realAnnotationPrefix string, patchAnnotationKey string) error {
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		panic(err)
	}
	// Creates the clientset
	clientSet, err := client.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// Don't use the downward API to fetch complicated annotations
	podAnnotations := getPodAnnotationsByName(clientSet, podNamespace, podName)
	patchPayload := getPatchPayload(podAnnotations, podName, realAnnotationPrefix, patchAnnotationKey)
	payloadBytes, _ := json.Marshal(patchPayload)
	result, err := clientSet.CoreV1().Pods(podNamespace).Patch(context.Background(), podName, types.JSONPatchType, payloadBytes, metav1.PatchOptions{}, "status")

	if err != nil {
		fmt.Println("Error:")
		fmt.Println(err)
	} else {
		fmt.Println("OK:")
		fmt.Println(result)
	}
	return err
}

func getPatchPayload(podAnnotations map[string]string, podName string, realAnnotationPrefix string, patchAnnotationKey string) []PatchOperation {
	var patch []PatchOperation

	annotationKey := fmt.Sprintf(realAnnotationPrefix + "/" + podName)
	valueMap := map[string]string{}

	for key, val := range podAnnotations {

		// Example: starting with this annotation: dd.replace/check_names: '["jmx"]'
		// Split into:
		// key=dd.replace/check_names
		// value='["jmx"]'

		// Only patch annotations with the patchAnnotationKey (dd.replace in this case)
		if !strings.Contains(key, patchAnnotationKey) {
			fmt.Println("Not matching with patchAnnKey:", key)
		}

		// Remove "dd.replace/" from key
		newKey := strings.Split(key, patchAnnotationKey+"/")[1]

		// Check values to see if we need to inject the full pod name
		newVal := strings.Replace(val, "trino-worker-", podName, -1)

		// Add the final values to the map
		valueMap[annotationKey+"."+newKey] = newVal
	}

	patch = append(patch, PatchOperation{
		Op:    "add",
		Path:  "/metadata/annotations",
		Value: valueMap,
	})
	return patch
}

func getPodAnnotationsByName(clientSet *client.Clientset, podNamespace string, podName string) map[string]string {
	pod, err := clientSet.CoreV1().Pods(podNamespace).Get(context.Background(), podName, metav1.GetOptions{})
	if err != nil {
		log.Fatal(err)
	}
	podAnnotations := pod.GetObjectMeta().GetAnnotations()
	return podAnnotations
}

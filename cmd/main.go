package main

import (
	"fmt"
	"github.com/rotemjac/k8s-patcher/pkg/patch"
	"io/ioutil"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"os"
)

func main() {
	// ############## Get ENV vars ############## //
	podInfoPath := os.Getenv("POD_INFO_PATH")
	podName := os.Getenv("POD_NAME")
	podNamespace := os.Getenv("POD_NAMESPACE")
	patchAnnotationKey := os.Getenv("PATCH_ANNOTATIONS_KEY")
	realAnnotationPrefix := os.Getenv("REAL_ANNOTATIONS_PREFIX")

	// ############## Fetch Labels ############## //
	podLabelsBytes, err := ioutil.ReadFile(podInfoPath + "/labels")
	if err != nil {
		fmt.Printf("Error reading podLabels file: %v\n", err)
		os.Exit(1)
	}
	podLabelString := string(podLabelsBytes)
	fmt.Println(podLabelString)

	// ############## Patch Annotations ############## //
	patch.Patch(podName, podNamespace, realAnnotationPrefix, patchAnnotationKey)
}

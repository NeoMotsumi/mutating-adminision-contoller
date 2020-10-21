package  handlers

import (
	"encoding/json"
	"net/http"

	"github.com/NeoMotsumi/mutating-adminision-contoller/pkg/utils"
	"github.com/NeoMotsumi/mutating-adminision-contoller/pkg/logger"
	"github.com/gorilla/mux"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
)

//RegisterMutatingWebhookHandlers registers all the webhook handlers.
func RegisterMutatingWebhookHandlers(r *mux.Router, lg logger.Logger)  {
	r.Handle("/mutate", mutatePod(lg))
}

type patchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}


//mutatePod is an http handler for all the pod mutation requests.
func mutatePod(lg logger.Logger) http.HandlerFunc {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {

		var pod corev1.Pod

		if err := utils.ReadRequest(req, &pod); err != nil {
			ar := reviewAdmission(&pod, lg)
		    utils.Respond(wr, http.StatusOK, ar)
		} else {
			utils.RespondErr(wr, err)
		}
	})
}

//reviewAdmission checks whether a pod requires a new label patch
func reviewAdmission(p *corev1.Pod, lg logger.Logger) *v1beta1.AdmissionResponse {

	if !mutationRequired(&p.ObjectMeta) {
		lg.Infof("Skipping mutation for %s/%s, label already appended.", p.Namespace, p.Name)
		return &v1beta1.AdmissionResponse{
			Allowed: true,
		}
	}

	patchBytes, err := createPatch(p)
	if err != nil {
		lg.Errorf(err.Error())
		return &v1beta1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}

	lg.Infof("AdmissionResponse: patch=%v\n", string(patchBytes))
	return &v1beta1.AdmissionResponse{
		Allowed: true,
		Patch:   patchBytes,
		PatchType: func() *v1beta1.PatchType {
			pt := v1beta1.PatchTypeJSONPatch
			return &pt
		}(),
	}
}

// createPatch creates a mutation patch for pod label
func createPatch(pod *corev1.Pod) ([]byte, error) {
	var patch []patchOperation

	patch = append(patch, patchOperation{
		Op:   "add",
		Path: "/metadata/labels",
		Value: map[string]string{
			"function": "webhook",
		},
	})

	return json.Marshal(patch)
}

//Validates whether the pod alread has a label for tier
func mutationRequired(metadata *metav1.ObjectMeta) bool {
	if _, ok := metadata.GetLabels()["function"]; ok {
		return false
	}

	return true
}

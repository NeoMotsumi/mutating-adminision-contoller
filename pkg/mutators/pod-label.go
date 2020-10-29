package mutators

import (
	"errors"
	"io/ioutil"
	"encoding/json"
	"net/http"

	v1beta1 "k8s.io/api/admission/v1beta1"
	"github.com/NeoMotsumi/mutating-adminision-contoller/pkg/utils"
	"github.com/NeoMotsumi/mutating-adminision-contoller/pkg/logger"
	 corev1 "k8s.io/api/core/v1"
	 metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	runtimeScheme = runtime.NewScheme()
	codecs        = serializer.NewCodecFactory(runtimeScheme)
	deserializer  = codecs.UniversalDeserializer()
)

type patchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

//MutatePod is an http handler for all the pod mutation requests.
func MutatePod(lg logger.Logger) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var body []byte
		if r.Body != nil {
			if data, err := ioutil.ReadAll(r.Body); err == nil {
				body = data
			}
		}

		ar := v1beta1.AdmissionReview{}
		if _, _, err := deserializer.Decode(body, nil, &ar); err != nil {
			lg.Errorf("Unable to decode body: %v", err)
			
			admissionErrResponse := &v1beta1.AdmissionReview{
				Response: &v1beta1.AdmissionResponse{
					Result: &metav1.Status{
						Message: err.Error(),
					},
				},
			}

			resp, err := json.Marshal(admissionErrResponse)
			if err != nil {
				lg.Errorf("Error Parsing admission review", err)
				utils.RespondErr(w, err, http.StatusBadRequest)
				return 
			}

			w.Write(resp)
			return
		}
		
		admissionResponse := reviewAdmission(&ar, lg)

		if admissionResponse != nil {
			ar.Response = admissionResponse
			ar.Response.UID = ar.Request.UID
			lg.Infof("AdmissionReview for Kind=%v, ApiVersion=%v",ar.Kind, ar.APIVersion)
		
			resp, err := json.Marshal(ar)

			if err != nil {
				lg.Errorf("Error Parsing review", err)
				utils.RespondErr(w, err, http.StatusInternalServerError)
				return 
			}

			if _, err := w.Write(resp); err != nil {
				lg.Errorf("Can't write response: %v", err)
				utils.RespondErr(w, err, http.StatusInternalServerError)
				return
			}
			
			lg.Infof("Webhook request handled successfully")
			return
		}

		err := errors.New("Could not process request")
		utils.RespondErr(w, err, http.StatusBadRequest)
	})
}

//reviewAdmission validates whether a pod requires a new label patch
func reviewAdmission(ar *v1beta1.AdmissionReview, lg logger.Logger) *v1beta1.AdmissionResponse {

	req := ar.Request
	var p corev1.Pod

	if err := json.Unmarshal(req.Object.Raw, &p); err != nil {
		lg.Errorf("Could not unmarshal raw object: %v", err)
		return &v1beta1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}

	lg.Infof("AdmissionReview for Kind=%v, Namespace=%v Name=%v (%v) UID=%v patchOperation=%v UserInfo=%v",
		req.Kind, req.Namespace, req.Name, p.Name, req.UID, req.Operation, req.UserInfo)

	if !mutationRequired(&p.ObjectMeta) {
		lg.Infof("Skipping mutation for %s/%s, label already appended.", p.Namespace, p.Name)
		return &v1beta1.AdmissionResponse{
			Allowed: true,
		}
	}

	patchBytes, err := createPatch(&p)
	if err != nil {
		lg.Errorf(err.Error())

		return &v1beta1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}

	r := &v1beta1.AdmissionResponse{ 
		Allowed: true,
		Patch:   patchBytes,
		PatchType: func() *v1beta1.PatchType {
			pt := v1beta1.PatchTypeJSONPatch
			return &pt
		}(),
	}

	lg.Infof("AdmissionReviewResponse for Kind=%v, ApiVersion=%v",ar.Kind, ar.APIVersion)
	return r
}

// createPatch creates a mutation patch for pod label
func createPatch(pod *corev1.Pod) ([]byte, error) {
	var patch []patchOperation

	lb := map[string]string{
		"function": "webhook",
	}

	for key, value := range lb {
		patch = append(patch, patchOperation{
			Op:   "add",
			Path: "/metadata/labels",
			Value: map[string]string{
				key: value,
			},
		})
	}

	return json.Marshal(patch)
}

//Validates whether the pod alread has a label for tier
func mutationRequired(metadata *metav1.ObjectMeta) bool {
	if _, ok := metadata.GetLabels()["function"]; ok {
		return false
	}

	return true
}

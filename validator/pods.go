package main

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
)

func podCheck(ns string, spec *v1.PodSpec) ([]string, error) {
	findings := []string{}

	for _, c := range spec.Containers {

		if c.ImagePullPolicy == v1.PullNever {
			findings = append(findings, fmt.Sprintf("Image %v exists? %v", c.Image, imageCheck(c.Image)))
		}

		if c.ImagePullPolicy == v1.PullAlways {
			findings = append(findings, fmt.Sprintf("Image %v pullable? %v", c.Image, imageCheck(c.Image)))
		}

		if c.ImagePullPolicy == v1.PullIfNotPresent {
			findings = append(findings, fmt.Sprintf("Image %v pullable? %v", c.Image, imageCheck(c.Image)))
		}

		// Referenced secret does not exists
		for _, env := range c.EnvFrom {
			if env.SecretRef != nil {
				secretName := env.SecretRef.Name
				res, err := checkSecret(secretName, ns, []string{})
				if err != nil && len(res) != 0 {
					return append(findings, res...), err
				}
				findings = append(findings, res...)
			}

			if env.ConfigMapRef != nil {
				configMapName := env.ConfigMapRef.Name
				res, err := checkConfigMap(configMapName, ns, []string{})
				if err != nil && len(res) != 0 {
					return append(findings, res...), err
				}
				findings = append(findings, res...)
			}
		}

		for _, env := range c.Env {
			if env.ValueFrom.ConfigMapKeyRef != nil {
				configMapName := env.ValueFrom.ConfigMapKeyRef.Name
				configMapKey := env.ValueFrom.ConfigMapKeyRef.Key
				res, err := checkConfigMap(configMapName, ns, []string{configMapKey})
				if err != nil && len(res) != 0 {
					return append(findings, res...), err
				}
				findings = append(findings, res...)
			}

			if env.ValueFrom.SecretKeyRef != nil {
				secretName := env.ValueFrom.SecretKeyRef.Name
				secretKey := env.ValueFrom.SecretKeyRef.Key
				res, err := checkSecret(secretName, ns, []string{secretKey})
				if err != nil && len(res) != 0 {
					return append(findings, res...), err
				}
				findings = append(findings, res...)
			}

		}
	}

	for _, v := range spec.Volumes {
		if v.Secret != nil {
			if v.Secret.SecretName != "" {
				res, err := checkSecret(v.Secret.SecretName, ns, []string{})
				if err != nil && len(res) != 0 {
					return append(findings, res...), err
				}
				findings = append(findings, res...)
			}
		}

		if v.ConfigMap != nil {
			if v.ConfigMap.Name != "" {
				keys := []string{}
				for _, ktp := range v.ConfigMap.Items {
					keys = append(keys, ktp.Key)
				}
				res, err := checkConfigMap(v.ConfigMap.Name, ns, []string{})
				if err != nil && len(res) != 0 {
					return append(findings, res...), err
				}
				findings = append(findings, res...)
			}
		}
	}

	return findings, nil
}

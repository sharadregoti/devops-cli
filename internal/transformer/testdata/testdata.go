package testdata

const PodJSON = `
{
    "apiVersion": "v1",
    "kind": "Pod",
    "metadata": {
        "annotations": {
            "kubectl.kubernetes.io/default-container": "httpbin",
            "kubectl.kubernetes.io/default-logs-container": "httpbin",
            "prometheus.io/path": "/stats/prometheus",
            "prometheus.io/port": "15020",
            "prometheus.io/scrape": "true",
            "sidecar.istio.io/status": "{\"initContainers\":[\"istio-init\"],\"containers\":[\"istio-proxy\"],\"volumes\":[\"istio-envoy\",\"istio-data\",\"istio-podinfo\",\"istio-token\",\"istiod-ca-cert\"],\"imagePullSecrets\":null,\"revision\":\"default\"}"
        },
        "creationTimestamp": "2022-07-18T12:07:34Z",
        "generateName": "httpbin-74fb669cc6-",
        "labels": {
            "app": "httpbin",
            "pod-template-hash": "74fb669cc6",
            "security.istio.io/tlsMode": "istio",
            "service.istio.io/canonical-name": "httpbin",
            "service.istio.io/canonical-revision": "v1",
            "version": "v1"
        },
        "name": "httpbin-74fb669cc6-8g9vf",
        "namespace": "default",
        "ownerReferences": [
            {
                "apiVersion": "apps/v1",
                "blockOwnerDeletion": true,
                "controller": true,
                "kind": "ReplicaSet",
                "name": "httpbin-74fb669cc6",
                "uid": "4981d2b6-b94e-4667-b20d-77227f17926b"
            }
        ],
        "resourceVersion": "6948789",
        "uid": "5bc1f218-baef-45a5-9367-99b7b65a2f90"
    },
    "spec": {
        "containers": [
            {
                "image": "docker.io/kennethreitz/httpbin",
                "imagePullPolicy": "IfNotPresent",
                "name": "httpbin",
                "ports": [
                    {
                        "containerPort": 80,
                        "protocol": "TCP"
                    }
                ],
                "resources": {},
                "terminationMessagePath": "/dev/termination-log",
                "terminationMessagePolicy": "File",
                "volumeMounts": [
                    {
                        "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount",
                        "name": "kube-api-access-v9csg",
                        "readOnly": true
                    }
                ]
            },
            {
                "args": [
                    "proxy",
                    "sidecar",
                    "--domain",
                    "$(POD_NAMESPACE).svc.cluster.local",
                    "--proxyLogLevel=warning",
                    "--proxyComponentLogLevel=misc:error",
                    "--log_output_level=default:info",
                    "--concurrency",
                    "2"
                ],
                "env": [
                    {
                        "name": "JWT_POLICY",
                        "value": "third-party-jwt"
                    },
                    {
                        "name": "PILOT_CERT_PROVIDER",
                        "value": "istiod"
                    },
                    {
                        "name": "CA_ADDR",
                        "value": "istiod.istio-system.svc:15012"
                    },
                    {
                        "name": "POD_NAME",
                        "valueFrom": {
                            "fieldRef": {
                                "apiVersion": "v1",
                                "fieldPath": "metadata.name"
                            }
                        }
                    },
                    {
                        "name": "POD_NAMESPACE",
                        "valueFrom": {
                            "fieldRef": {
                                "apiVersion": "v1",
                                "fieldPath": "metadata.namespace"
                            }
                        }
                    },
                    {
                        "name": "INSTANCE_IP",
                        "valueFrom": {
                            "fieldRef": {
                                "apiVersion": "v1",
                                "fieldPath": "status.podIP"
                            }
                        }
                    },
                    {
                        "name": "SERVICE_ACCOUNT",
                        "valueFrom": {
                            "fieldRef": {
                                "apiVersion": "v1",
                                "fieldPath": "spec.serviceAccountName"
                            }
                        }
                    },
                    {
                        "name": "HOST_IP",
                        "valueFrom": {
                            "fieldRef": {
                                "apiVersion": "v1",
                                "fieldPath": "status.hostIP"
                            }
                        }
                    },
                    {
                        "name": "PROXY_CONFIG",
                        "value": "{}\n"
                    },
                    {
                        "name": "ISTIO_META_POD_PORTS",
                        "value": "[\n    {\"containerPort\":80,\"protocol\":\"TCP\"}\n]"
                    },
                    {
                        "name": "ISTIO_META_APP_CONTAINERS",
                        "value": "httpbin"
                    },
                    {
                        "name": "ISTIO_META_CLUSTER_ID",
                        "value": "Kubernetes"
                    },
                    {
                        "name": "ISTIO_META_INTERCEPTION_MODE",
                        "value": "REDIRECT"
                    },
                    {
                        "name": "ISTIO_META_WORKLOAD_NAME",
                        "value": "httpbin"
                    },
                    {
                        "name": "ISTIO_META_OWNER",
                        "value": "kubernetes://apis/apps/v1/namespaces/default/deployments/httpbin"
                    },
                    {
                        "name": "ISTIO_META_MESH_ID",
                        "value": "cluster.local"
                    },
                    {
                        "name": "TRUST_DOMAIN",
                        "value": "cluster.local"
                    }
                ],
                "image": "docker.io/istio/proxyv2:1.12.1",
                "imagePullPolicy": "IfNotPresent",
                "name": "istio-proxy",
                "ports": [
                    {
                        "containerPort": 15090,
                        "name": "http-envoy-prom",
                        "protocol": "TCP"
                    }
                ],
                "readinessProbe": {
                    "failureThreshold": 30,
                    "httpGet": {
                        "path": "/healthz/ready",
                        "port": 15021,
                        "scheme": "HTTP"
                    },
                    "initialDelaySeconds": 1,
                    "periodSeconds": 2,
                    "successThreshold": 1,
                    "timeoutSeconds": 3
                },
                "resources": {
                    "limits": {
                        "cpu": "2",
                        "memory": "1Gi"
                    },
                    "requests": {
                        "cpu": "10m",
                        "memory": "40Mi"
                    }
                },
                "securityContext": {
                    "allowPrivilegeEscalation": false,
                    "capabilities": {
                        "drop": [
                            "ALL"
                        ]
                    },
                    "privileged": false,
                    "readOnlyRootFilesystem": true,
                    "runAsGroup": 1337,
                    "runAsNonRoot": true,
                    "runAsUser": 1337
                },
                "terminationMessagePath": "/dev/termination-log",
                "terminationMessagePolicy": "File",
                "volumeMounts": [
                    {
                        "mountPath": "/var/run/secrets/istio",
                        "name": "istiod-ca-cert"
                    },
                    {
                        "mountPath": "/var/lib/istio/data",
                        "name": "istio-data"
                    },
                    {
                        "mountPath": "/etc/istio/proxy",
                        "name": "istio-envoy"
                    },
                    {
                        "mountPath": "/var/run/secrets/tokens",
                        "name": "istio-token"
                    },
                    {
                        "mountPath": "/etc/istio/pod",
                        "name": "istio-podinfo"
                    },
                    {
                        "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount",
                        "name": "kube-api-access-v9csg",
                        "readOnly": true
                    }
                ]
            }
        ],
        "dnsPolicy": "ClusterFirst",
        "enableServiceLinks": true,
        "initContainers": [
            {
                "args": [
                    "istio-iptables",
                    "-p",
                    "15001",
                    "-z",
                    "15006",
                    "-u",
                    "1337",
                    "-m",
                    "REDIRECT",
                    "-i",
                    "*",
                    "-x",
                    "",
                    "-b",
                    "*",
                    "-d",
                    "15090,15021,15020"
                ],
                "image": "docker.io/istio/proxyv2:1.12.1",
                "imagePullPolicy": "IfNotPresent",
                "name": "istio-init",
                "resources": {
                    "limits": {
                        "cpu": "2",
                        "memory": "1Gi"
                    },
                    "requests": {
                        "cpu": "10m",
                        "memory": "40Mi"
                    }
                },
                "securityContext": {
                    "allowPrivilegeEscalation": false,
                    "capabilities": {
                        "add": [
                            "NET_ADMIN",
                            "NET_RAW"
                        ],
                        "drop": [
                            "ALL"
                        ]
                    },
                    "privileged": false,
                    "readOnlyRootFilesystem": false,
                    "runAsGroup": 0,
                    "runAsNonRoot": false,
                    "runAsUser": 0
                },
                "terminationMessagePath": "/dev/termination-log",
                "terminationMessagePolicy": "File",
                "volumeMounts": [
                    {
                        "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount",
                        "name": "kube-api-access-v9csg",
                        "readOnly": true
                    }
                ]
            }
        ],
        "nodeName": "docker-desktop",
        "preemptionPolicy": "PreemptLowerPriority",
        "priority": 0,
        "restartPolicy": "Always",
        "schedulerName": "default-scheduler",
        "securityContext": {},
        "serviceAccount": "httpbin",
        "serviceAccountName": "httpbin",
        "terminationGracePeriodSeconds": 30,
        "tolerations": [
            {
                "effect": "NoExecute",
                "key": "node.kubernetes.io/not-ready",
                "operator": "Exists",
                "tolerationSeconds": 300
            },
            {
                "effect": "NoExecute",
                "key": "node.kubernetes.io/unreachable",
                "operator": "Exists",
                "tolerationSeconds": 300
            }
        ],
        "volumes": [
            {
                "emptyDir": {
                    "medium": "Memory"
                },
                "name": "istio-envoy"
            },
            {
                "emptyDir": {},
                "name": "istio-data"
            },
            {
                "downwardAPI": {
                    "defaultMode": 420,
                    "items": [
                        {
                            "fieldRef": {
                                "apiVersion": "v1",
                                "fieldPath": "metadata.labels"
                            },
                            "path": "labels"
                        },
                        {
                            "fieldRef": {
                                "apiVersion": "v1",
                                "fieldPath": "metadata.annotations"
                            },
                            "path": "annotations"
                        }
                    ]
                },
                "name": "istio-podinfo"
            },
            {
                "name": "istio-token",
                "projected": {
                    "defaultMode": 420,
                    "sources": [
                        {
                            "serviceAccountToken": {
                                "audience": "istio-ca",
                                "expirationSeconds": 43200,
                                "path": "istio-token"
                            }
                        }
                    ]
                }
            },
            {
                "configMap": {
                    "defaultMode": 420,
                    "name": "istio-ca-root-cert"
                },
                "name": "istiod-ca-cert"
            },
            {
                "name": "kube-api-access-v9csg",
                "projected": {
                    "defaultMode": 420,
                    "sources": [
                        {
                            "serviceAccountToken": {
                                "expirationSeconds": 3607,
                                "path": "token"
                            }
                        },
                        {
                            "configMap": {
                                "items": [
                                    {
                                        "key": "ca.crt",
                                        "path": "ca.crt"
                                    }
                                ],
                                "name": "kube-root-ca.crt"
                            }
                        },
                        {
                            "downwardAPI": {
                                "items": [
                                    {
                                        "fieldRef": {
                                            "apiVersion": "v1",
                                            "fieldPath": "metadata.namespace"
                                        },
                                        "path": "namespace"
                                    }
                                ]
                            }
                        }
                    ]
                }
            }
        ]
    },
    "status": {
        "conditions": [
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2022-12-10T02:40:17Z",
                "status": "True",
                "type": "Initialized"
            },
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2022-12-10T02:41:16Z",
                "status": "True",
                "type": "Ready"
            },
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2022-12-10T02:41:16Z",
                "status": "True",
                "type": "ContainersReady"
            },
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2022-07-18T12:07:34Z",
                "status": "True",
                "type": "PodScheduled"
            }
        ],
        "containerStatuses": [
            {
                "containerID": "docker://0a146c8f7db2ba687bf1194396b23912eacb769399d83046b9bf4736a83ea658",
                "image": "kennethreitz/httpbin:latest",
                "imageID": "docker-pullable://kennethreitz/httpbin@sha256:599fe5e5073102dbb0ee3dbb65f049dab44fa9fc251f6835c9990f8fb196a72b",
                "lastState": {
                    "terminated": {
                        "containerID": "docker://ab410538551b5676a917c403bad3f0364a291519c80d282b77022c0eff3b3576",
                        "exitCode": 255,
                        "finishedAt": "2022-12-10T02:39:20Z",
                        "reason": "Error",
                        "startedAt": "2022-12-08T06:11:18Z"
                    }
                },
                "name": "httpbin",
                "ready": true,
                "restartCount": 33,
                "started": true,
                "state": {
                    "running": {
                        "startedAt": "2022-12-10T02:40:17Z"
                    }
                }
            },
            {
                "containerID": "docker://8d97de73702965588dbafd5f6395ecb5346408d2f9a4026ff0c991ccf68f0b6b",
                "image": "istio/proxyv2:1.12.1",
                "imageID": "docker-pullable://istio/proxyv2@sha256:4704f04f399ae24d99e65170d1846dc83d7973f186656a03ba70d47bd1aba88f",
                "lastState": {
                    "terminated": {
                        "containerID": "docker://cc0515c90a2fdfef74f2e5160686a98ddc6d63ff68a4d18c9eb52fab9d4a459e",
                        "exitCode": 255,
                        "finishedAt": "2022-12-10T02:39:20Z",
                        "reason": "Error",
                        "startedAt": "2022-12-08T06:11:18Z"
                    }
                },
                "name": "istio-proxy",
                "ready": true,
                "restartCount": 33,
                "started": true,
                "state": {
                    "running": {
                        "startedAt": "2022-12-10T02:40:17Z"
                    }
                }
            }
        ],
        "hostIP": "192.168.65.4",
        "initContainerStatuses": [
            {
                "containerID": "docker://242152cb4b439a1a398616c290fad3f33b34816f18402f14c39eaf664a84ea62",
                "image": "istio/proxyv2:1.12.1",
                "imageID": "docker-pullable://istio/proxyv2@sha256:4704f04f399ae24d99e65170d1846dc83d7973f186656a03ba70d47bd1aba88f",
                "lastState": {},
                "name": "istio-init",
                "ready": true,
                "restartCount": 33,
                "state": {
                    "terminated": {
                        "containerID": "docker://242152cb4b439a1a398616c290fad3f33b34816f18402f14c39eaf664a84ea62",
                        "exitCode": 0,
                        "finishedAt": "2022-12-10T02:40:15Z",
                        "reason": "Completed",
                        "startedAt": "2022-12-10T02:40:15Z"
                    }
                }
            }
        ],
        "phase": "Running",
        "podIP": "10.1.17.12",
        "podIPs": [
            {
                "ip": "10.1.17.12"
            }
        ],
        "qosClass": "Burstable",
        "startTime": "2022-07-18T12:07:34Z"
    }
}`

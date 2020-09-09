// Code generated by go generate; DO NOT EDIT.
// using data from templates/kubernetes
package kubernetes

func TemplatesMap() map[string]string {
	templatesMap := make(map[string]string)

	templatesMap["cluster_role.yaml"] = `apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  labels:
    app: cf-argocd-agent
  name: cf-argocd-agent
rules:
  - apiGroups:
      - argoproj.io
    resources:
      - applications
      - appprojects
    verbs:
      - get
      - list
      - watch
`

	templatesMap["cluster_role_binding.yaml"] = `apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  labels:
    app: cf-argocd-agent
  name: cf-argocd-agent
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cf-argocd-agent
subjects:
  - kind: ServiceAccount
    name: cf-argocd-agent
    namespace: {{ .Namespace }}`

	templatesMap["deployment.yaml"] = `apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: cf-argocd-agent
  name: cf-argocd-agent
  namespace: {{ .Namespace }}
spec:
  selector:
    matchLabels:
      app: cf-argocd-agent
  replicas: 1
  revisionHistoryLimit: 5
  strategy:
    rollingUpdate:
      maxSurge: 50%
      maxUnavailable: 50%
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: cf-argocd-agent
    spec:
      serviceAccountName: cf-argocd-agent
      containers:
      - env:
        - name: ARGO_HOST
          value: {{ .Argo.Host }}
        - name: ARGO_USERNAME
          value: {{ .Argo.Username }}
        - name: ARGO_PASSWORD
          value: {{ .Argo.Password }}
        - name: CODEFRESH_HOST
          value: {{ .Codefresh.Host }}
        - name: CODEFRESH_TOKEN
          value: {{ .Codefresh.Token }}
        - name: IN_CLUSTER
          value: "true"
        - name: CODEFRESH_INTEGRATION
          value: {{ .Codefresh.Integration }}
        image: codefresh/argocd-agent:stable
        imagePullPolicy: Always
        name: cf-argocd-agent
      restartPolicy: Always`

	templatesMap["sa.yaml"] = `apiVersion: v1
kind: ServiceAccount
metadata:
  labels:
    app: cf-argocd-agent
  name: cf-argocd-agent
  namespace: {{ .Namespace }}`

	return templatesMap
}

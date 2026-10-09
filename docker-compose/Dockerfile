FROM quay.io/argoproj/argocd:v3.4.3@sha256:b6d161b4984c4bffb8c02250252fdd17f23fa4e2aa25fada0e606e1fa08d2348 AS argocd
FROM zegl/kube-score:v1.20.0@sha256:ac4c43ad560af905d66f6bf57b0937c591332e6dbf2167c31193a13b4695ab97 as kube-score
FROM ghcr.io/yannh/kubeconform:v0.8.0@sha256:faffaf43f95aa6425306e1ab8d6fcad72acb9049158f38e574c085ea1ec0f64e as kubeconform
FROM ghcr.io/kyverno/kyverno-cli:v1.19.1@sha256:ced7b2be0b04250cabfe695f15307f69eb715fe23234816388af4f3812915b2a AS kyverno

FROM ghcr.io/coopnorge/engineering-docker-images/e0/devtools-golang-v1beta1:latest@sha256:627ccb7fec3a9666ceb8d85014aaa0e5da2f9530428009cc427cfc1c1ba2f8df AS golang-devtools

COPY --from=kube-score /kube-score /usr/local/bin/kube-score
COPY --from=argocd /usr/local/bin/argocd /usr/local/bin/argocd
COPY --from=argocd /usr/local/bin/helm /usr/local/bin/helm
COPY --from=argocd /usr/local/bin/kustomize /usr/local/bin/kustomize
COPY --from=kubeconform /kubeconform /usr/local/bin/kubeconform
COPY --from=kyverno /ko-app/kubectl-kyverno /usr/local/bin/kyverno

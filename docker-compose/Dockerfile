FROM quay.io/argoproj/argocd:v3.4.3@sha256:b6d161b4984c4bffb8c02250252fdd17f23fa4e2aa25fada0e606e1fa08d2348 AS argocd
FROM zegl/kube-score:v1.20.0@sha256:ac4c43ad560af905d66f6bf57b0937c591332e6dbf2167c31193a13b4695ab97 as kube-score
FROM ghcr.io/yannh/kubeconform:v0.8.0@sha256:faffaf43f95aa6425306e1ab8d6fcad72acb9049158f38e574c085ea1ec0f64e as kubeconform
FROM ghcr.io/kyverno/kyverno-cli:v1.18.1@sha256:b7e272572d244ddec0b83469f7200ba883555bf69de4b294cee52a197c8c6590 AS kyverno

FROM ghcr.io/coopnorge/engineering-docker-images/e0/devtools-golang-v1beta1:latest@sha256:b7759b7beedd07bbb59e22bf278291879b17257f34afc5755bd6d8ffad7bcc2d AS golang-devtools

COPY --from=kube-score /kube-score /usr/local/bin/kube-score
COPY --from=argocd /usr/local/bin/argocd /usr/local/bin/argocd
COPY --from=argocd /usr/local/bin/helm /usr/local/bin/helm
COPY --from=argocd /usr/local/bin/kustomize /usr/local/bin/kustomize
COPY --from=kubeconform /kubeconform /usr/local/bin/kubeconform
COPY --from=kyverno /ko-app/kubectl-kyverno /usr/local/bin/kyverno

FROM quay.io/argoproj/argocd:v3.3.8@sha256:5d45dc6db21db32a0638ac9128462c6d9956a90fc81760146dada5a243ff7516 AS argocd
FROM zegl/kube-score:v1.20.0@sha256:ac4c43ad560af905d66f6bf57b0937c591332e6dbf2167c31193a13b4695ab97 as kube-score
FROM ghcr.io/yannh/kubeconform:v0.7.0@sha256:85dbef6b4b312b99133decc9c6fc9495e9fc5f92293d4ff3b7e1b30f5611823c as kubeconform
FROM ghcr.io/kyverno/kyverno-cli:v1.17.1@sha256:3a7addf83897abae291361dfccbdf63cba5926c7795fe44e527883ae2237e4b8 AS kyverno

FROM ghcr.io/coopnorge/engineering-docker-images/e0/devtools-golang-v1beta1:latest@sha256:399d1d97b15a414f3d46ad90c0584cbb409d9f14629daec9a67e1955615b3435 AS golang-devtools

COPY --from=kube-score /kube-score /usr/local/bin/kube-score
COPY --from=argocd /usr/local/bin/argocd /usr/local/bin/argocd
COPY --from=argocd /usr/local/bin/helm /usr/local/bin/helm
COPY --from=argocd /usr/local/bin/kustomize /usr/local/bin/kustomize
COPY --from=kubeconform /kubeconform /usr/local/bin/kubeconform
COPY --from=kyverno /ko-app/kubectl-kyverno /usr/local/bin/kyverno

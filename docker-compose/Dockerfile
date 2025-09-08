FROM quay.io/argoproj/argocd:v3.1.4@sha256:dd6239528d15e407a2e1ad20fe77771a3aacfa4680701e00454dce58230f15eb AS argocd
FROM zegl/kube-score:v1.20.0@sha256:ac4c43ad560af905d66f6bf57b0937c591332e6dbf2167c31193a13b4695ab97 as kube-score
FROM ghcr.io/yannh/kubeconform:v0.7.0@sha256:85dbef6b4b312b99133decc9c6fc9495e9fc5f92293d4ff3b7e1b30f5611823c as kubeconform
FROM ghcr.io/kyverno/kyverno-cli:v1.15.1@sha256:6e40920087f32e24948dd12d94f47eb4cb4d6d3e110b564b79510ed78a9bb8a6 AS kyverno

FROM ghcr.io/coopnorge/engineering-docker-images/e0/devtools-golang-v1beta1:latest@sha256:a3df482d6c4ee86a34699d3f51ed88c931d9b2f678730160f8a5170cd4f8c2a8 AS golang-devtools

COPY --from=kube-score /kube-score /usr/local/bin/kube-score
COPY --from=argocd /usr/local/bin/argocd /usr/local/bin/argocd
COPY --from=argocd /usr/local/bin/helm /usr/local/bin/helm
COPY --from=argocd /usr/local/bin/kustomize /usr/local/bin/kustomize
COPY --from=kubeconform /kubeconform /usr/local/bin/kubeconform
COPY --from=kyverno /ko-app/kubectl-kyverno /usr/local/bin/kyverno

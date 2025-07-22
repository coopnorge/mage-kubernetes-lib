FROM quay.io/argoproj/argocd:v3.0.11@sha256:1cb4ede2fc4a6332c128d9ba29d19e8cb5b76f9260dc47550a4b3b154269ff86 AS argocd
FROM zegl/kube-score:v1.20.0@sha256:ac4c43ad560af905d66f6bf57b0937c591332e6dbf2167c31193a13b4695ab97 as kube-score
FROM ghcr.io/yannh/kubeconform:v0.7.0@sha256:85dbef6b4b312b99133decc9c6fc9495e9fc5f92293d4ff3b7e1b30f5611823c as kubeconform
FROM ghcr.io/kyverno/kyverno-cli:v1.14.4@sha256:bcc4db143edb795cb6a20b111b60368925a9b2ebdc66284afe748c3f7c27dab4 AS kyverno

FROM ghcr.io/coopnorge/engineering-docker-images/e0/devtools-golang-v1beta1:latest@sha256:e1594c4147e3af785d1e8ef7bb0423c6b1cf4e14ac57fb3ef3fdf803bee3a739 AS golang-devtools

COPY --from=kube-score /kube-score /usr/local/bin/kube-score
COPY --from=argocd /usr/local/bin/argocd /usr/local/bin/argocd
COPY --from=argocd /usr/local/bin/helm /usr/local/bin/helm
COPY --from=argocd /usr/local/bin/kustomize /usr/local/bin/kustomize
COPY --from=kubeconform /kubeconform /usr/local/bin/kubeconform
COPY --from=kyverno /ko-app/kubectl-kyverno /usr/local/bin/kyverno

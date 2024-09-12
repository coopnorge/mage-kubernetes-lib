
FROM docker.io/argoproj/argocd:v2.6.15@sha256:58ebb4ed23c8db4cd4cc3f954f8d94c4b4e3d9669c0751c484108d22b86d52de as argocd
FROM zegl/kube-score:v1.17.0@sha256:d4c34b9560d9ad35a5113fbbff6d798c87e6986a306ee097a0eb4242f9714810 as kube-score
FROM ghcr.io/yannh/kubeconform:v0.6.4@sha256:e68a0b638c6e9b76f1b7d58b4ec94340ef3b6601db25b2e40b29e3ac2d68e4bf as kubeconform

FROM ghcr.io/coopnorge/engineering-docker-images/e0/devtools-golang-v1beta1:latest@sha256:146ecad927da14e3ea5ea36800f2b8b3117ddc488f9343c38f8c544092347332 AS golang-devtools

COPY --from=kube-score /kube-score /usr/local/bin/kube-score
COPY --from=argocd /usr/local/bin/argocd /usr/local/bin/argocd
COPY --from=argocd /usr/local/bin/helm /usr/local/bin/helm
COPY --from=argocd /usr/local/bin/kustomize /usr/local/bin/kustomize
COPY --from=kubeconform /kubeconform /usr/local/bin/kubeconform

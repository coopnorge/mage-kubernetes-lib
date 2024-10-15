
FROM docker.io/argoproj/argocd:v2.6.15@sha256:58ebb4ed23c8db4cd4cc3f954f8d94c4b4e3d9669c0751c484108d22b86d52de as argocd
FROM zegl/kube-score:v1.19.0@sha256:94137f32ce139dc9fbdbbd380249025e4d378c282ff151a100b981cdeeb923b6 as kube-score
FROM ghcr.io/yannh/kubeconform:v0.6.7@sha256:0925177fb05b44ce18574076141b5c3d83235e1904d3f952182ac99ddc45762c as kubeconform

FROM ghcr.io/coopnorge/engineering-docker-images/e0/devtools-golang-v1beta1:latest@sha256:c295378915f6ab5e2db0bdc4138642bc391028f8b1d4c85e2510046727733262 AS golang-devtools

COPY --from=kube-score /kube-score /usr/local/bin/kube-score
COPY --from=argocd /usr/local/bin/argocd /usr/local/bin/argocd
COPY --from=argocd /usr/local/bin/helm /usr/local/bin/helm
COPY --from=argocd /usr/local/bin/kustomize /usr/local/bin/kustomize
COPY --from=kubeconform /kubeconform /usr/local/bin/kubeconform

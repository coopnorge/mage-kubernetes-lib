
FROM docker.io/argoproj/argocd:v2.6.15@sha256:58ebb4ed23c8db4cd4cc3f954f8d94c4b4e3d9669c0751c484108d22b86d52de as argocd
FROM zegl/kube-score:v1.19.0@sha256:94137f32ce139dc9fbdbbd380249025e4d378c282ff151a100b981cdeeb923b6 as kube-score
FROM ghcr.io/yannh/kubeconform:v0.6.7@sha256:0925177fb05b44ce18574076141b5c3d83235e1904d3f952182ac99ddc45762c as kubeconform
FROM ghcr.io/kyverno/kyverno-cli:v1.14.1@sha256:591dcdafc5d5e4ba4d9a5f8d59975af73fd4a8f93fdab6398aacec0f0f65783e AS kyverno

FROM ghcr.io/coopnorge/engineering-docker-images/e0/devtools-golang-v1beta1:latest@sha256:3ec237d4993c3f558cd7eb71439dd57428e855e877304c34595f3e80adceea22 AS golang-devtools

COPY --from=kube-score /kube-score /usr/local/bin/kube-score
COPY --from=argocd /usr/local/bin/argocd /usr/local/bin/argocd
COPY --from=argocd /usr/local/bin/helm /usr/local/bin/helm
COPY --from=argocd /usr/local/bin/kustomize /usr/local/bin/kustomize
COPY --from=kubeconform /kubeconform /usr/local/bin/kubeconform
COPY --from=kyverno /ko-app/kubectl-kyverno /usr/local/bin/kyverno
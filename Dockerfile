ARG IMAGE_TOOLS
ARG IMAGE_MICRO
ARG IMAGE_BASE

FROM $IMAGE_TOOLS as tools
ARG CGO_ENABLED
ENV CGO_ENABLED ${CGO_ENABLED}
ENV GOCACHE /.cache/go-build
ENV GOMODCACHE /.cache/mod
# HACK: Otherwise the build is not cached, no idea why
RUN touch /ok
WORKDIR /src

FROM tools as kubectl
ARG KUBECTL_URL
RUN curl -sLf "$KUBECTL_URL" > /usr/bin/kubectl && \
    chmod +x /usr/bin/kubectl

FROM tools as loglevel
ARG LOGLEVEL_URL
RUN curl -sLf "$LOGLEVEL_URL" | tar xvzf - -C /usr/bin

FROM tools as controllergen
ARG CONTROLLER_GEN_VERSION
RUN go install sigs.k8s.io/controller-tools/cmd/controller-gen@${CONTROLLER_GEN_VERSION}

FROM tools as build-go
COPY ./go.mod ./go.sum ./
COPY ./pkg/apis/go.mod ./pkg/apis/go.sum ./pkg/apis/
COPY ./pkg/client/go.mod ./pkg/client/go.sum ./pkg/client/
WORKDIR /src/pkg/apis
RUN --mount=type=cache,target=/.cache go mod download
WORKDIR /src/pkg/client
RUN --mount=type=cache,target=/.cache go mod download
WORKDIR /src
RUN --mount=type=cache,target=/.cache go mod download
RUN --mount=type=cache,target=/.cache go mod verify

FROM tools as build-go-codegen
COPY ./ ./
ARG GO_LDFLAGS
ARG GO_GCFLAGS
ARG GO_BUILDFLAGS
RUN --mount=type=cache,target=/.cache go build ${GO_BUILDFLAGS} -gcflags "${GO_GCFLAGS}" -ldflags "${GO_LDFLAGS}" -o /main ./pkg/codegen/
RUN --mount=type=cache,target=/.cache go build ${GO_BUILDFLAGS} -gcflags "${GO_GCFLAGS}" -ldflags "${GO_LDFLAGS}" -o /buildconfig ./pkg/codegen/buildconfig/
RUN --mount=type=cache,target=/.cache go build ${GO_BUILDFLAGS} -gcflags "${GO_GCFLAGS}" -ldflags "${GO_LDFLAGS}" -o /cleanup ./pkg/codegen/generator/cleanup/

# TODO: controller-gen
FROM build-go as build-go-generator
COPY --link --from=controllergen /go/bin/controller-gen /usr/bin/
COPY --link --from=build-go-codegen /main /usr/bin/
COPY --link --from=build-go-codegen /buildconfig /usr/bin/
COPY --link --from=build-go-codegen /cleanup /usr/bin/
COPY ./scripts/ ./scripts/
COPY ./pkg/apis/ ./pkg/apis/
COPY ./pkg/buildconfig/ ./pkg/buildconfig/
COPY ./*.go ./build.yaml ./
RUN --mount=type=cache,target=/.cache go generate ./...

FROM scratch as go-generator
COPY --link --from=build-go-generator /src/pkg/apis /pkg/apis/
COPY --link --from=build-go-generator /src/pkg/generated /pkg/generated/

FROM build-go as build-server
COPY ./ ./
COPY --link --from=go-generator / ./
ARG GO_LDFLAGS
ARG GO_GCFLAGS
ARG GO_BUILDFLAGS
RUN --mount=type=cache,target=/.cache go build ${GO_BUILDFLAGS} -gcflags "${GO_GCFLAGS}" -ldflags "${GO_LDFLAGS}" -o /rancher

FROM build-go as build-agent
COPY ./ ./
ARG GO_LDFLAGS
ARG GO_GCFLAGS
ARG GO_BUILDFLAGS
RUN --mount=type=cache,target=/.cache go build ${GO_BUILDFLAGS} -gcflags "${GO_GCFLAGS}" -ldflags "${GO_LDFLAGS}" -o /agent ./cmd/agent

FROM scratch as binary-server
COPY --link --from=build-server /rancher /rancher

FROM scratch as binary-agent
COPY --link --from=build-agent /agent /agent

FROM $IMAGE_BASE as agent
COPY --link --from=kubectl /usr/bin/kubectl /usr/bin/kubectl
COPY --link --from=loglevel /usr/bin/ /usr/bin/
COPY --link --from=binary-agent /agent /usr/bin/agent

# RUN go build -tags k8s \
#   -gcflags="all=${GCFLAGS}" \
#   -ldflags \
#   "-X github.com/rancher/rancher/pkg/version.Version=$VERSION
#    -X github.com/rancher/rancher/pkg/version.GitCommit=$COMMIT
#    -X github.com/rancher/rancher/pkg/settings.InjectDefaults=$DEFAULT_VALUES $LINKFLAGS" \
#   -o bin/rancher
#

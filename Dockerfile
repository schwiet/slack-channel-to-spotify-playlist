FROM golang:1.18 as build

WORKDIR /usr/src/handler

# cache dependencies
COPY go.mod go.sum ./
RUN go mod download & go mod verify

# build
COPY lambda .
ARG FUNC_PATH
RUN go build -o handler $FUNC_PATH

# copy artifacts to a clean image
FROM public.ecr.aws/lambda/go:1
COPY --from=build /usr/src/handler ${LAMBDA_TASK_ROOT}
CMD [ "handler" ]

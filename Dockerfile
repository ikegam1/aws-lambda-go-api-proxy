FROM public.ecr.aws/bitnami/golang:1.18 as build-image

WORKDIR /go/src
#COPY go.mod go.sum main.go ./
Add . ./

#RUN go mod init github.com/awslabs/aws-lambda-go-api-proxy \
#  && go mod edit -replace github.com/awslabs/aws-lambda-go-api-proxy=github.com/ikegam1/aws-lambda-go-api-proxy@issue-144 \
#  && go get -u github.com/ikegam1/aws-lambda-go-api-proxy@issue-144 \
#  && go mod tidy

#RUN go clean -modcache
#RUN go get -u github.com/ikegam1/aws-lambda-go-api-proxy@issue-144 
#RUN go get -u github.com/aws/aws-lambda-go
RUN go get -u github.com/labstack/echo/v4@v4.7.2 
RUN go get github.com/labstack/echo/v4/middleware@v4.7.2

RUN go build
RUN ls -lha ./

FROM public.ecr.aws/lambda/go:1

COPY --from=build-image /go/src/ /var/task/

# Command can be overwritten by providing a different command in the template directly.
CMD ["aws-lambda-go-api-proxy"]

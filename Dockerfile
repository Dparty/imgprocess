FROM golang:1.17 AS build-stage
WORKDIR /app
COPY . .
RUN go build -o main

FROM golang:1.17 AS production-stage
WORKDIR /app
COPY --from=build-stage /app/default_img default_img/
COPY --from=build-stage /app/main ./main
EXPOSE 80
CMD ["./main", "80", "https://ordering-uat-1318552943.cos.ap-hongkong.myqcloud.com"]
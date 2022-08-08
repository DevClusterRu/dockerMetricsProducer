FROM 355982936287.dkr.ecr.us-east-1.amazonaws.com/golang1.18-alpine:latest
WORKDIR /app
COPY . .
RUN go mod tidy
CMD ["go", "run", "main.go"]

#EXPOSE 4470
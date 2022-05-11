FROM golang:1.17.10-alpine3.15

RUN apk add gcc tesseract-ocr-dev gcc libc-dev npm git --no-cache
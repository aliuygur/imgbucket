FROM scratch
ADD imgbucket /app/imgbucket
EXPOSE 5000
ENTRYPOINT ["/app/imgbucket"]
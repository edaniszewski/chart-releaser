FROM scratch

LABEL org.label-schema.schema-version="1.0" \
      org.label-schema.name="chartreleaser/chart-releaser" \
      org.label-schema.vcs-url="https://github.com/edaniszewski/chart-releaser" \
      org.label-schema.vendor="Erick Daniszewski"

ADD chart-releaser /bin

WORKDIR /release

ENTRYPOINT ["/bin/chart-releaser"]

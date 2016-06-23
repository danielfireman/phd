# This image is based on the popular Alpine Linux project, available in the
# alpine official image. Alpine Linux is much smaller than most distribution
# base images (~5MB), and thus leads to much slimmer images in general.
FROM java:8-jre-alpine

COPY ./target/restserver-1.0.jar /usr/bin/
WORKDIR /usr/bin
CMD ["java",  "-jar", "restserver-1.0.jar"]

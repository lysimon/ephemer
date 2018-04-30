FROM jenkins:latest

MAINTAINER laurent.simon

# Copying default configuration to disable jenkins security
COPY config.xml /var/jenkins_home/config.xml

# Install jenkins plugin
RUN /usr/local/bin/install-plugins.sh git
RUN /usr/local/bin/install-plugins.sh s3


# Disable first time initialization
ENV JAVA_OPTS -Dhudson.Main.development=true -Djenkins.install.runSetupWizard=false


# Install VI for easier maintenance
#RUN apt-get update
#RUN apt-get install -y vim sudo
USER root

FROM jenkins:latest

MAINTAINER laurent.simon

USER root
RUN apt-get update \
      && apt-get install -y sudo \
      && rm -rf /var/lib/apt/lists/*


# Copying default configuration to disable jenkins security
COPY config.xml /var/jenkins_home/config.xml

# Install jenkins pipeline plugin
RUN /usr/local/bin/install-plugins.sh workflow-aggregator


# Disable first time initialization
ENV JAVA_OPTS -Dhudson.Main.development=true -Djenkins.install.runSetupWizard=false


# Install VI for easier maintenance
#RUN apt-get update
#RUN apt-get install -y vim sudo

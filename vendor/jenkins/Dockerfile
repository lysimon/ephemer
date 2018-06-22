FROM jenkins:latest

MAINTAINER laurent.simon

USER root
RUN apt-get update \
      && apt-get install -y sudo \
      && rm -rf /var/lib/apt/lists/*


# Copying default configuration to disable jenkins security
COPY config.xml /var/jenkins_home/config.xml
# Copying our default jobs to jenkins
COPY jobs /var/jenkins_home/jobs

# Install jenkins pipeline plugin
RUN /usr/local/bin/install-plugins.sh git


# Disable first time initialization
ENV JAVA_OPTS -Dhudson.Main.development=true -Djenkins.install.runSetupWizard=false

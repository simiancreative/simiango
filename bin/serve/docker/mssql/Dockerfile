FROM mcr.microsoft.com/mssql/server

# Switch to root user for access to apt-get install
USER root

# Install node/npm
RUN apt-get -y update  && \
        apt-get install -y dos2unix

# Create app directory
RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

# Bundle app source
COPY . /usr/src/app

RUN dos2unix *

RUN chmod +x /usr/src/app/import-data.sh

# Switch back to mssql user and run the entrypoint script
USER mssql

ENV MSSQL_SA_PASSWORD=Melon40twomonkeys
ENV ACCEPT_EULA=true

RUN /usr/src/app/import-data.sh

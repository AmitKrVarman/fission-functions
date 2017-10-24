FISSION Go Boilerplate
=========================


Getting Started
-----------------
 - Install Go 1.9 (https://golang.org/dl/)
 - Ensure your GOPATH is set such that $PWD == $GOPATH/src/github.com/AmitKrVarman/fission-functions
 - TBC
 - TBC
 - TBC
 - TBC

Go-boilerplate features
-----------------------

The Fission go-boilerplate is a simple template that --


Exercising the boilerplate handlers
-----------------------------------
``` shell

```
And the corresponding log output:
``` shell

```
Steps
1. Install docker.io - Login with your Docker ID to push and pull images from Docker Hub. If you don't have a Docker ID, head over to https://hub.docker.com to create one.

2. Install Go Lang -> $sudo apt-get install golang-go
3. Stern https://github.com/wercker/stern
4. gcloud - https://cloud.google.com/sdk/downloads#apt-get

project ID is landg-179815, 
gcloud container clusters list
gcloud container clusters get-credentials my-first-cluster

5. concourse
https://concourse.landg.madeden.net
 
username: concourse
password: concourse
 
download the fly CLI on your laptop, and set it up to start configuring the pipelines. First things to look at are nodejs and go. Documentation is on https://concourse.ci
 
6. 

Errors

E1 - fission function create --name process-form --env go --deploy function.so
Failed to upload file function.so: invalid character '<' looking for beginning of value

sudo docker login
curl https://index.docker.io/v1/

{"errors":[{"code":"UNAUTHORIZED","message":"authentication required","detail":null}]}

docker run -p 5000:5000 registry

apt-get install --only-upgrade docker-engine

dpkg was interrupted, you must manually run 'sudo dpkg --configure -a' to correct the problem.

sudo apt-get update
$ sudo apt-get upgrade

Network Issues

E2 Stern Erros
gcloud config set container/use_client_certificate True 
export CLOUDSDK_CONTAINER_USE_CLIENT_CERTIFICATE=True
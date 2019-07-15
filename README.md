# Usage

Automate the clone of all gitlab projects with to OpenGrok.
maps groups and subgroups of projects in a tree.

## Dependencies:

Use go get to install some needed libraries:

```go
go get "github.com/xanzy/go-gitlab"
```
```go
go get "gopkg.in/src-d/go-git.v4"
```

## Build
Run this command in in the path where **opengrok-gitlab-sync** lives.
```go
go build
```

## Usage
The build step should generate an executable, Run it using the following parameters **in order**:
- Group id (vistaprint-org group id).
- Token (Should have read access to all projects in the organization).
- Path (provide the path where you want to checkout the projects).

## Deployment
This work is deployed on GKE

## Authentication/SSO
Integration with Auth0 is done using oauth2-proxy on kubernetes
Implementation details: https://vistaprint.atlassian.net/wiki/spaces/TNT/pages/33161217/Oauth2+proxy+implementation+with+Auth0+in+Kubernetes

## Certification management:
The certificate is issued by Let's Encrypt for opengrok.vips.vistaprint.io
The automation of the creation and renewal of the certificate is managed by cert-manager within Kubernetes
Implementation details: https://vistaprint.atlassian.net/wiki/spaces/TNT/pages/32899085/Certificates+management+on+Kubernetes

## To be done:
IaC, automate the rolling updates on Kubernetes from Gitlab, helm chart, requests and limits for pods, ResourceQuota for opengrok namespace


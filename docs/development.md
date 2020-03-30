# Development

## Setup Tools

- CRC

  Follow CRC setup guidelines to setup your development environment. You can find the documentation [here](https://cloud.redhat.com/openshift/install/crc/installer-provisioned).

- OpenShift CLI (oc)

-  ko. Get the installation steps [here](https://github.com/google/ko).

## Deploy Application on CRC


### Deploy API Service

Ensure you are in  `cd backend/api` directory

#### Prepare API Release Yaml
Export `KO_DOCKER_REPO` for ko to publish image to. E.g.

```
export KO_DOCKER_REPO=quay.io/<username>
```

`ko` resolve and apply the `api.yaml`

```
cd backend/api
ko resolve -f config > api.yaml
```

The command above will create a container image and push it to the registry
pointed by `KO_DOCKER_REPO`. Ensure that the image is **publicly** available

#### Update the GitHub Api secret and token

Edit `api.yaml` and update the secret - `api`. Set GitHub `oauth` client id and
secret, access token.

```
apiVersion: v1
kind: Secret
metadata:
  name: api
  namespace: tekton-hub
type: Opaque
stringData:
  GITHUB_TOKEN: My Personal access token
  CLIENT_ID: Oauth client id
  CLIENT_SECRET: Oauth secret
```

**NOTE:** DO NOT MODIFY `config/20-api-secret.yaml` commit and push


#### Apply API Release Yaml

```
oc apply -f api.yaml
```

Watch the pods until `db` is running. `api` pod will fail at this stage as
`db` is not created yet.

```
oc get pods -o wide -w
```

At this stage the `deployement` `db` should be up and running.

#### Create Database

Ensure `db` pod is `running`

```
$ oc get pods

NAME                   READY   STATUS    RESTARTS   AGE
api-6675fbf9f5-fft4h   0/1     Error     3          72s
db-748f56cb8c-rwqjc    1/1     Running   1          72s
  ^^^^^^^^^^^^^               ^^^^^^^^^
  <db-pod-name>
```

Create database `tekton_hub` by `rsh` in the `db` pod.  You can get the pod name from above command or use `oc get pod -l app=db`

```
oc rsh <db-pod-name>
```

Once you are in the pod, use `psql` to create the database

```
psql  -c 'create database tekton_hub;'
```

#### Ensure api service is running

At this stage, `api` should be in `Running` state

```
$ oc get pods

NAME                   READY   STATUS    RESTARTS   AGE
api-6675fbf9f5-fft4h   0/1     Running   3          72s
                               ^^^^^^^
db-748f56cb8c-rwqjc    1/1     Running   1          72s

```

Now, Both the pods are up but the database is empty.

To create tables and initialise the data, we need to run the db-migration.

Run the below command to create migration image and the job.

```
ko apply -f config/db-migration/14-db-migration.yaml
```

Check the logs using ` oc logs job/db-migration `.

Wait till the migration log shows
```
Migration did run successfully !!
```

#### Verify if api route is accessible

```
curl -k -X GET -I $(oc get routes api --template='https://{{ .spec.host }}/resources')
```

### Deploy Validation Service

Generating validation release yaml follows the same instructions as the [Api
Service][] above


```
cd backend/validation
```

1. Export `KO_DOCKER_REPO` for ko to publish image
1. `ko` resolve and apply the `validation.yaml`
  ```
    cd backend/validation
    ko resolve -f config > validation.yaml
  ```
1. Ensure the `validation` image is public
1. `oc apply -f validation.yaml`
1. Watch for valiation pod to be running


### Deploy UI

```
cd frontend
```

#### Build and Publish Image

```
docker build -t <image> . && docker push <image>
```
#### Update the deployment image

Update `config/11-deployement` to use the image built above

#### Update the GitHub OAuth Client ID

Edit `config/10-config.yaml` and set your GitHub OAuth Client ID

```
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: ui
  namespace: tekton-hub
data:
  API_URL: 'https://api-tekton-hub.apps-crc.testing'
  GH_CLIENT_ID: GH OAuth Client ID   <<< update this
```

#### Apply the manifests

```
oc apply -f config
```

#### Ensure pods are up and running

```
oc get pods -o wide -w
```

Open: oc get routes ui --template='https://{{ .spec.host }} '

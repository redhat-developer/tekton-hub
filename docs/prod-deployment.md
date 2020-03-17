# Deploy to Production

## Prerequisites

- OpenShift Cluster
- OpenShift CLI (oc)
- ko. You can find installation steps [here](https://github.com/google/ko).

## Deployment

- [Deploy from Scratch](#Deploy-from-Scratch)
- [Upgrade to a new release](#Upgrade-to-a-new-release)

## Deploy from Scratch

### Step 1: Create a git tag

Create a new git tag and use the same tag to tag the images. You can see the previous tag [here](https://github.com/redhat-developer/tekton-hub/releases).

### Step 2: Deploy API Service and Database

Ensure you are in `cd backend/api` directory.
```
cd backend/api
```

#### Prepare API Release Yaml
Export `KO_DOCKER_REPO` for ko to publish image to quay.io/tekton-hub.

```
export KO_DOCKER_REPO=quay.io/tekton-hub
```

`ko` resolve and apply the `api.yaml`

Use git tag created in the Step 1 and replace it with `<git-tag>` to tag the image in the below command. eg. `-t v0.3`.

Make sure you are logged in to the quay.io/tekton-hub.

```
ko resolve -t <git-tag> -f config/scratch/ > api.yaml
```

The command above will create a container image and push it to the `quay.io/tekton-hub`.

#### Update the GitHub Api secret, token and Image name

Edit `api.yaml` and update the secret - `api`. Set GitHub `oauth` client id and secret, access token.

```
apiVersion: v1
kind: Secret
metadata:
  name: api
  namespace: tekton-hub
type: Opaque
stringData:
  GITHUB_TOKEN: My Personal access token       <<<
  CLIENT_ID: Oauth Client Id                   <<< Update this values
  CLIENT_SECRET: Oauth Secret                  <<<
```

Update the `POSTGRESQL_PASSWORD` in `db` secret. Use random password for db.
```
apiVersion: v1
kind: Secret
metadata:
  name: db
  namespace: tekton-hub
type: Opaque
stringData:
  POSTGRESQL_DATABASE: tekton_hub
  POSTGRESQL_USER: postgres
  POSTGRESQL_PASSWORD: Database Password   <<<  Update this value
  POSTGRESQL_PORT: "5432"
```

Update the image name in `api.yaml` to look like as below. Remove the sha from image name.

```
spec:
      containers:
      - name: api
        image: quay.io/tekton-hub/api-b786b59ff17bae6:v0.3   <<<  Update here
        ports:
        - containerPort: 5000
```

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
                              ^^^^^^^^^
```

Connect to database by port-forwarding `db` service

```
oc port-forward svc/db 5432:5432
```

On a different terminal, use `psql` to create and load the database

```
psql -h localhost -U postgres -p 5432 -c 'create database tekton_hub;'
psql -h localhost -U postgres -p 5432 tekton_hub < backups/02-01-2020.dump
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
>NOTE: you may want to end the port-forward session.

#### Verify if api route is accessible

```
curl -k -X GET -I $(oc get routes api --template='https://{{ .spec.host }}/resources')
```

### Step 3: Deploy Validation Service

Generating validation release yaml follows the same instructions as the Api
Service above.


```
cd backend/validation
```

Export `KO_DOCKER_REPO` for ko to publish image to quay.io/tekton-hub.

```
export KO_DOCKER_REPO=quay.io/tekton-hub
```
`ko` resolve and apply the `validation.yaml`.

Use git tag created in Step 1 and replace it with `<git-tag>` to tag the image in below command. eg. `-t v0.3`
```
ko resolve -t <git-tag> -f config/scratch/ > validation.yaml
```
Update the image name in validation.yaml to look like as below. Remove the `sha` from image name.


```
...
spec:
      containers:
      - name: api
        image: quay.io/tekton-hub/validation-b786b59ff17bae6:v0.3    <<< Update here
        ports:
        - containerPort: 5000
...
```


#### Apply Validation Release Yaml
```
oc apply -f validation.yaml
```


### Step 4: Deploy UI

```
cd frontend
```

#### Build and Publish Image

Use image name as `ui` and git tag from Step 1 and replace it with `<git-tag>` for tagging the image in below command. eg. `-t ui:v03`

Make sure you are logged in to the quay.io/tekton-hub.
```
docker build -t quay.io/tekton-hub/ui:<git-tag> . && docker push quay.io/tekton-hub/ui:<git-tag>
```
#### Update the deployment image

Update `config/scratch/11-deployement` to use the image built above.
```
...
 containers:
        - name: ui
          image: quay.io/tekton-hub/ui:v03     <<< Update Image Name with tag
...
```

#### Update the GitHub OAuth Client ID

Edit `config/scratch/10-config.yaml` and set your GitHub OAuth Client ID and Api Service Route as `API_URL`.

You can use `oc get routes api --template='https://{{ .spec.host }}'` to get the Api service route.

```
apiVersion: v1
kind: ConfigMap
metadata:
  name: ui
  namespace: tekton-hub
data:
  API_URL: Api Service route             <<<   Update this values
  GH_CLIENT_ID: GH OAuth Client ID       <<<
```

#### Apply the manifests

```
oc apply -f config/scratch/
```

#### Ensure pods are up and running

```
oc get pods -o wide -w
```

### Tekton Hub url

```
oc get routes ui --template='https://{{ .spec.host }}
```



## Upgrade to a new release

### Step 1: Create a git tag

Create a new git tag and use the same tag to tag the images. You can see the previous tag [here](https://github.com/redhat-developer/tekton-hub/releases).

### Step 2: Deploy API Service

Ensure you are in `cd backend/api` directory.
```
cd backend/api
```

#### Prepare API Release Yaml
Export `KO_DOCKER_REPO` for ko to publish image to quay.io/tekton-hub.

```
export KO_DOCKER_REPO=quay.io/tekton-hub
```

`ko` resolve and apply the `api.yaml`

Use git tag created in the Step 1 and replace it with `<git-tag>` to tag the image in the below command. eg. `-t v0.3`.

Make sure you are logged in to the quay.io/tekton-hub.

```
ko resolve -t <git-tag> -f config/upgrade/ > api.yaml
```

The command above will create a container image and push it to the `quay.io/tekton-hub`.

#### Update the Image name

Update the image name in `api.yaml` to look like as below. Remove the sha from image name.

```
spec:
      containers:
      - name: api
        image: quay.io/tekton-hub/api-b786b59ff17bae6:v0.3   <<<  Update here
        ports:
        - containerPort: 5000
```

#### Apply API Release Yaml

```
oc apply -f api.yaml
```

Watch the pods until new `api` pod is running.
```
oc get pods -o wide -w
```

### Step 3: Deploy Validation Service

Generating validation release yaml follows the same instructions as the Api
Service above.

```
cd backend/validation
```

Export `KO_DOCKER_REPO` for ko to publish image to quay.io/tekton-hub.

```
export KO_DOCKER_REPO=quay.io/tekton-hub
```
`ko` resolve and apply the `validation.yaml`.

Use git tag created in Step 1 and replace it with `<git-tag>` to tag the image in below command. eg. `-t v0.3`
```
ko resolve -t <git-tag> -f config/upgrade/ > validation.yaml
```
Update the image name in validation.yaml to look like as below. Remove the `sha` from image name.


```
...
spec:
      containers:
      - name: api
        image: quay.io/tekton-hub/validation-b786b59ff17bae6:v0.3    <<< Update here
        ports:
        - containerPort: 5000
...
```


#### Apply Validation Release Yaml
```
oc apply -f validation.yaml
```

#### Ensure new pod is up and running

```
oc get pods -o wide -w
```


### Step 4: Deploy UI

```
cd frontend
```

#### Build and Publish Image

Use image name as `ui` and git tag from Step 1 and replace it with `<git-tag>` for tagging the image in below command. eg. `-t ui:v03`

Make sure you are logged in to the quay.io/tekton-hub.
```
docker build -t quay.io/tekton-hub/ui:<git-tag> . && docker push quay.io/tekton-hub/ui:<git-tag>
```
#### Update the deployment image

Update `config/upgrade/11-deployement` to use the image built above.
```
...
 containers:
        - name: ui
          image: quay.io/tekton-hub/ui:v03     <<< Update Image Name with tag
...
```

#### Apply the deployment manifest

```
oc apply -f config/upgrade/
```

#### Ensure new pod is up and running

```
oc get pods -o wide -w
```

### Tekton Hub url

```
oc get routes ui --template='https://{{ .spec.host }}
```
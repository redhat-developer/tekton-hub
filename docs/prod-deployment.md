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

Create a new git tag and use the same tag to tag the images. You can see the
previous tag [here](https://github.com/redhat-developer/tekton-hub/releases).

Checkout the latest code and create a new tag.
```
git tag -as <tag-name>
```
And push it to upstream/master.
```
git push upstream/master <tagname>
```

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

Use git tag created in the Step 1 and replace it with `<git-tag>` to tag the
image in the below command. eg. `-t v0.3`.

Make sure you are logged in to the quay.io/tekton-hub.

```
ko resolve -t <git-tag> -f config/ > api.yaml
```

The command above will create a container image and push it to the
`quay.io/tekton-hub`.

#### Update the GitHub Api secret, token and Image name

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

Update the image name in `api.yaml` to look like as below. Remove the sha from
image name.

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
   ^^^^^^^^^^^^               ^^^^^^^^^
   db-pod-name
```

Create database `tekton_hub` by rsh into the db pod. You can get the pod name
from above command or use `oc get pod -l app=db`
```
oc rsh <db-pod-name>
```
Once you are in the pod, use `psql` to create the database.
```
psql -c 'create database tekton_hub;'
```

Use `exit` to get out off the pod.


#### Ensure api service is running

At this stage, `api` should be in `Running` state

```
$ oc get pods

NAME                   READY   STATUS    RESTARTS   AGE
api-6675fbf9f5-fft4h   0/1     Running   3          72s
                               ^^^^^^^
db-748f56cb8c-rwqjc    1/1     Running   1          72s

```
If the api pod is still not in running state, try deleting the pod, new pod
will be created.

```
oc delete pod <api-pod-name>
```
Now, Both the pods are up but the database is empty.

To create tables and initialise the data, we need to run the db-migration.

Run the below command to create migration image.

Use git tag created in the Step 1 and replace it with <git-tag> to tag the image in the below command. eg. `-t v0.3`.

```
ko resolve -f config/db-migration -t <git-tag>
```

The Database migration should be ran only once. So, we will run a kubernetes
job.

Edit the `config/db-migration/14-db-migration.yaml` and update the image name
to look like as below. Remove the `sha` from image name.

```
...
 spec:
      containers:
      - name: db-migration
        image: quay.io/tekton-hub/db-e1225b1694ead695:v0.4   <<< Update here
...
```

Apply the migration job yaml.

```
oc apply -f config/db-migration/14-db-migration.yaml
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

Use git tag created in Step 1 and replace it with `<git-tag>` to tag the image
in below command. eg. `-t v0.3`
```
ko resolve -t <git-tag> -f config/ > validation.yaml
```
Update the image name in validation.yaml to look like as below. Remove the
`sha` from image name.


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

Use image name as `ui` and git tag from Step 1 and replace it with `<git-tag>`
for tagging the image in below command. eg. `ui:v03`

Make sure you are logged in to the quay.io/tekton-hub.
```
docker build -t quay.io/tekton-hub/ui:<git-tag> . &&
  docker push quay.io/tekton-hub/ui:<git-tag>
```
#### Update the deployment image

Update `config/11-deployement` to use the image built above.
```
...
 containers:
        - name: ui
          image: quay.io/tekton-hub/ui:v03     <<< Update Image Name with tag
...
```

#### Update the GitHub OAuth Client ID

Edit `config/10-config.yaml` and set your GitHub OAuth Client ID and Api
Service Route as `API_URL`.

You can use `oc get routes api --template='https://{{ .spec.host }}'` to get
the Api service route.

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
oc apply -f config/
```

#### Ensure pods are up and running

```
oc get pods -o wide -w
```

### Tekton Hub url

```
oc get routes ui --template='https://{{ .spec.host }}'
```

## Upgrade to a new release

### Step 1: Create a git tag


Create a new git tag and use the same tag to tag the images. You can see the previous tags [here](https://github.com/redhat-developer/tekton-hub/releases).

Checkout the latest code and create a new tag.
```
git tag <tag-name>
```
And push it to upstream/master.
```
git push upstream/master <tagname>
```

### Step 2: Run DB Migration

Run the below command to create migration image.

Use git tag created in the Step 1 and replace it with <git-tag> to tag the image in the below command. eg. `-t v0.3`.

```
ko resolve -f config/db-migration -t <git-tag>
```

The Database migration should be ran only once. So, we will run a kubernetes job.

Edit the `config/db-migration/14-db-migration.yaml` and update the image name to look like as below. Remove the `sha` from image name.

```
...
 spec:
      containers:
      - name: db-migration
        image: quay.io/tekton-hub/db-e1225b1694ead695:v0.3   <<< Update here
...
```

Apply the migration job yaml.

```
oc apply -f config/db-migration/14-db-migration.yaml
```

Check the logs using ` oc logs job/db-migration `.

Wait till the migration log shows
```
Migration did run successfully !!
```

### Step 3: Deploy API Service

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
ko resolve -t <git-tag> -f config/21-api-deployment.yaml > api.yaml
```

The command above will create a container image and push it to the `quay.io/tekton-hub`.

#### Update the Image name

Update the image name in `api.yaml` to look like as below. Remove the `sha` from image name.

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
ko resolve -t <git-tag> -f config/10-deployment.yaml > validation.yaml
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


### Step 5: Deploy UI

```
cd frontend
```

#### Build and Publish Image

Use image name as `ui` and git tag from Step 1 and replace it with `<git-tag>` for tagging the image in below command. eg. `ui:v03`

Make sure you are logged in to the quay.io/tekton-hub.
```
docker build -t quay.io/tekton-hub/ui:<git-tag> . && docker push quay.io/tekton-hub/ui:<git-tag>
```
#### Update the deployment image

Update `config/11-deployement` to use the image built above.
```
...
 containers:
        - name: ui
          image: quay.io/tekton-hub/ui:v03     <<< Update Image Name with tag
...
```

#### Apply the deployment manifest

```
oc apply -f config/11-deployement
```

#### Ensure new pod is up and running

```
oc get pods -o wide -w
```

### Tekton Hub url

```
oc get routes ui --template='https://{{ .spec.host }}
```

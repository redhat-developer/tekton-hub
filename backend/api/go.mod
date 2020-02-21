module github.com/redhat-developer/tekton-hub/backend/api

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/ghodss/yaml v1.0.0
	github.com/google/go-cmp v0.3.1 // indirect
	github.com/google/go-github v17.0.0+incompatible
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/gorilla/handlers v1.4.2
	github.com/gorilla/mux v1.7.3
	github.com/joho/godotenv v1.3.0
	github.com/lib/pq v1.2.0
	github.com/mattbaird/jsonpatch v0.0.0-20171005235357-81af80346b1a // indirect
	github.com/tektoncd/pipeline v0.9.2
	go.uber.org/zap v1.13.0 // indirect
	golang.org/x/oauth2 v0.0.0-20191202225959-858c2ad4c8b6
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543 // indirect
	gopkg.in/yaml.v2 v2.2.7 // indirect
	k8s.io/api v0.17.0 // indirect
	k8s.io/client-go v0.0.0-20190805141520-2fe0317bcee0 // indirect
	k8s.io/utils v0.0.0-20191218082557-f07c713de883 // indirect
	knative.dev/pkg v0.0.0-20191230183737-ead56ad1f3bd // indirect
)

go 1.13

# validation-service

This validation-service do lint and schema validation for tasks and pipelines during uploading task/pipeline on Tekton Hub.
### Dependencies
1. Go 1.11.3

### Running on your local machine 
1. Fork and clone this repository
```
git clone https://github.com/redhat-developer/tekton-hub
cd tekton-hub/backend/validation

```
2. Install dependencies 
``` go mod download ```
3. Build the application
```go build main.go```
4. Run the application
``` ./main```



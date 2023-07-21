# helm-shortlist


### Notes
---

- https://helm.sh/docs/topics/advanced/
- https://pkg.go.dev/helm.sh/helm/v3@v3.12.2/pkg/action#Configuration
- https://cobra.dev/
- https://github.com/helm/helm/blob/main/cmd/helm/list_test.go
- https://dev.to/kcdchennai/how-to-create-your-first-helm-plugin-4i0g


#### bitnami helm charts 
```bash

$ helm repo add bitnami https://charts.bitnami.com/bitnami

$ helm install bitnami/mysql --generate-name
```

####  deploy simple helm charts with pod manifest only 

```bash

$ helm install simple-helm-chart simple_helm_chart/ 

NAME: test
LAST DEPLOYED: Fri Jul 14 23:14:17 2023
NAMESPACE: default
STATUS: deployed
REVISION: 1
TEST SUITE: None
```

first output from the script 

```
mifot@ThinkPad-X1-Carbon-4th:~/GoProjects/helm-shortlist$ go run list.go 
2023/07/14 23:15:05 &{Name:test Info:0xc000328000 Chart:0xc0000d9a40 Config:map[] Manifest:---
# Source: simple-chart/templates/pod.yaml
apiVersion: v1
kind: Pod
metadata:
  name: simple-pod
spec:
  containers:
    - name: simple-chart
      image: "nginx:1.25.1"
      imagePullPolicy: IfNotPresent
 Hooks:[] Version:1 Namespace:default Labels:map[modifiedAt:1689369258 name:test owner:helm status:deployed version:1]}
```

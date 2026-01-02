# Start cloud provider running in our cluster - just the standard k3s components:
# CoreDNS, local-path-provisioner, and metrics-server. 
# This is a clean slate.
colima start --arch x86_64 --cpu 2 --memory 3 --disk 30 --mount-type sshfs --runtime docker --dns 8.8.8.8 --dns 1.1.1.1 --kubernetes --network-address

# Cluster.....................................................................
# if `kubectl cluster-info` not exists:
kind create cluster --name calendar-k8s --config kind-cluster-config.yaml

# check cluster
kubectl config get-contexts
CURRENT   NAME                CLUSTER             AUTHINFO            NAMESPACE
          colima              colima              colima              linkding
          default                                                     
          ida-colima          linkding            linkding            
*         kind-calendar-k8s   kind-calendar-k8s   kind-calendar-k8s   
          linkding

kubectl cluster-info

To further debug and diagnose cluster problems, use 'kubectl cluster-info dump'.

# Ingress ......................................................................
# An Ingress controller receives external HTTP(S) traffic and routes it to Services inside
# the cluster according to Ingress resources
# if OK
kubectl get ingressclass
NAME    CONTROLLER             PARAMETERS   AGE
nginx   k8s.io/ingress-nginx   <none>       117s

# else:
helm install ingress-nginx ingress-nginx/ingress-nginx \
  --namespace ingress-nginx \
  --create-namespace

kubectl get pods -n ingress-nginx
NAME READY   STATUS    RESTARTS   AGE
ingress-nginx-controller-6b48bd57b6-nwfbw   1/1     Running   0          78s

kubectl get validatingwebhookconfiguration
NAME                      WEBHOOKS   AGE
ingress-nginx-admission   1          64s


# App stage ................................................................
helm upgrade --install calendar-app . \
  --namespace calendar \
  --create-namespace

# make sure which Docker context is used:
docker context ls
NAME            DESCRIPTION                               DOCKER ENDPOINT                                 ERROR
colima *        colima                                    unix:///Users/mdx/.colima/default/docker.sock   
default         Current DOCKER_HOST based configuration   unix:///var/run/docker.sock                     
desktop-linux                                             unix:///Users/mdx/.docker/run/docker.sock


 
# Check k3s resources of the app:
kubectl get pods -n calendar


kubectl get svc -n calendar
NAME           TYPE       CLUSTER-IP    EXTERNAL-IP   PORT(S)          AGE
calendar-app   NodePort   10.96.94.17   <none>        8081:30090/TCP   4m11s
 

kubectl logs calendar-app-75ddd987cb-hbl56 --tail=20 -n calendar
Defaulted container "app" out of: app, wait-for-postgres (init)
POSTGRES_DSN is set. Running database migrations...
2025/12/26 15:27:04 goose: no migrations to run. current version: 1
Migrations complete.
Starting application...
HTTP gateway listening on :8081
gRPC server listening on :50051
gRPC call start: /calendarGRPC.CalendarService/HealthCheck
health check requested
gRPC call end: /calendarGRPC.CalendarService/HealthCheck | duration: 21.701µs

curl http://myapp.local/health
{"status":"OK"}

kubectl exec -it deploy/calendar -- curl http://producer:8082
kubectl exec -it deploy/calendar -- curl http://consumer:8083


# App stop and cleanup ................................................................
helm uninstall calendar-app -n calendar
kind delete cluster --name calendar-k8s



kubectl logs ingress-nginx-controller-5fd9b9dddd-49pn9 -n ingress-nginx


kubectl get pods -n calendar -o name \
| xargs -I {} kubectl logs -n calendar {} \
  --all-containers=true \
  --prefix=true \
  --tail=100
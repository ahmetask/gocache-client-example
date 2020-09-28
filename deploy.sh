minikube start
eval $(minikube docker-env)
docker build -t gocache-client .
kubectl delete deployment gocache-client
kubectl delete service gocache-client
kubectl create -f deployment.yml
kubectl create -f service.yml
minikube service gocache-client

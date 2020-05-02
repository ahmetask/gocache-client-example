minikube start
eval $(minikube docker-env)
docker build -t gocache-client .
kubectl delete deployment gocache-client
kubectl delete service gocache-client
kubectl apply -f deployment.yml
kubectl apply -f service.yml
minikube service gocache-client
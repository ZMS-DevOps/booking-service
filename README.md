# Example 

Turn on ingress 

```shell
minikube addons enable ingress
```

Create and delete namespace:

```bash
kubectl create namespace backend
kubectl delete namespace backend
```
Sve kubernetes fajlove pokrenuti da dobijemo configmap, secret i mongo service i statefulSet
```shell
kubectl -n backend apply -f mongo-configmap.yml
kubectl -n backend apply -f mongo-secret.yml
kubectl -n backend apply -f search-configmap.yml 
kubectl -n backend apply -f mongo.yml
kubectl apply -f search-service.yml 
```
Get pods:
```shell
kubectl -n backend get pods
kubectl get pods
```

Testing load balancing and service:
```shell
kubectl -n backend run -it --rm  --image curlimages/curl:8.00.1 curl -- sh
kubectl run -it --rm  --image curlimages/curl:8.00.1 curl -- sh
```
Inside the container execute `curl http://hotel:8083/hotels` (hotel jer je to naziv servisa)
```shell
 curl http://booking:8086/booking/health


RESERVATION REQUESTS
curl http://booking:8086/booking/unavailability


curl http://booking:8086/booking/unavailability
curl http://booking/booking/unavailability -H "Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJSM19wNUs5dFVtTkpmSVlhdmo1NV8xTjI1VUJOMmh5cmJJWGVnVmV4ZG1ZIn0.eyJleHAiOjE3MTYwMzM2MTgsImlhdCI6MTcxNjAzMzMxOCwianRpIjoiNGFkMmZmM2ItYTAwYS00NmQxLWE4ZDktNDNjNjJmNWUyMTA3IiwiaXNzIjoiaHR0cDovL2tleWNsb2FrLmRlZmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWw6ODA4MC9yZWFsbXMvSXN0aW8iLCJhdWQiOiJhY2NvdW50Iiwic3ViIjoiNTY5YjkxMjktN2NjYi00OWNhLTgzY2YtMDZmNmMzMGZmM2NjIiwidHlwIjoiQmVhcmVyIiwiYXpwIjoiSXN0aW8iLCJzZXNzaW9uX3N0YXRlIjoiMzkyMWU0NWMtOWNkNS00ZWE5LWFjZTktMTZhM2MyMzUzMzdiIiwiYWNyIjoiMSIsImFsbG93ZWQtb3JpZ2lucyI6WyIvKiJdLCJyZWFsbV9hY2Nlc3MiOnsicm9sZXMiOlsib2ZmbGluZV9hY2Nlc3MiLCJhZG1pbiIsInVtYV9hdXRob3JpemF0aW9uIiwiZGVmYXVsdC1yb2xlcy1pc3RpbyJdfSwicmVzb3VyY2VfYWNjZXNzIjp7ImFjY291bnQiOnsicm9sZXMiOlsibWFuYWdlLWFjY291bnQiLCJtYW5hZ2UtYWNjb3VudC1saW5rcyIsInZpZXctcHJvZmlsZSJdfX0sInNjb3BlIjoicHJvZmlsZSBlbWFpbCIsInNpZCI6IjM5MjFlNDVjLTljZDUtNGVhOS1hY2U5LTE2YTNjMjM1MzM3YiIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJuYW1lIjoiTWlsYW4gQWpkZXIiLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJob3RlbC1hZG1pbiIsImdpdmVuX25hbWUiOiJNaWxhbiIsImZhbWlseV9uYW1lIjoiQWpkZXIiLCJlbWFpbCI6ImFqZGVyLm1pbGFuMjAwMEBnbWFpbC5jb20ifQ.JBPiqSMdgcuuhDu8tJ3RAm7Vq0TVZD0KrcuKvObRXKuC7XsKAfQVX0mvd7JZYsmjgIdOMHYC6f8HxC-sXbdXvgzAJeCyx9icW0hfYQ1qfuizW_AfFnTbmqovZFeJME6c2n_okUKM8bNXOHWNvVR0uZ6UdrBiyNemNfXRjSKA8SAMiIOB_ZbH4aUqEtrNbZ3CEe3G8ZMUJ8UlYdzUeBt2VxPZenBqwCkhmYAx4pmdkuF5f0pB8sMKsbpjgkssgYnLwalg9tDRmxLNOv7eeBenyiZm73EwFfnE0xvf0KtmauwYBUkx91AFZfuvlRBM1u4O8PGfo24iMV_S524GfbSMsA" 
 
curl http://booking/booking/unavailability/664b8a267cf6a66bcf5fc587
 
curl --location --request PUT 'http://booking:8086/booking/unavailability/add' \
--header 'Content-Type: application/json' \
--data '{
    "accommodation_id": "6643a56c9dea1760db469b7b",
    "start": "2025-04-18T00:00:00Z",
    "end": "2025-04-28T00:00:00Z",
    "reason": "OwnerSet"
}'

curl --location --request PUT 'http://booking/booking/unavailability/remove' \
--header 'Content-Type: application/json' \
--data '{
    "unavailability_id": "6649ed2589ddc2e5a3d852a9",
    "start": "2025-04-19T00:00:00Z",
    "end": "2025-04-21T00:00:00Z",
    "reason": "OwnerSet"
}'

Update price
curl -X PUT http://hotel/hotel/accommodation/price/66467094ea31a50a941a503e -d '{"date_range": {"start": "2024-05-15T00:00:00Z", "end": "2024-05-20T00:00:00Z"}, "price": 100.50, "type": "PerGuest"}' -H "Content-Type: application/json" -H "Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJSM19wNUs5dFVtTkpmSVlhdmo1NV8xTjI1VUJOMmh5cmJJWGVnVmV4ZG1ZIn0.eyJleHAiOjE3MTYwMzM2MTgsImlhdCI6MTcxNjAzMzMxOCwianRpIjoiNGFkMmZmM2ItYTAwYS00NmQxLWE4ZDktNDNjNjJmNWUyMTA3IiwiaXNzIjoiaHR0cDovL2tleWNsb2FrLmRlZmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWw6ODA4MC9yZWFsbXMvSXN0aW8iLCJhdWQiOiJhY2NvdW50Iiwic3ViIjoiNTY5YjkxMjktN2NjYi00OWNhLTgzY2YtMDZmNmMzMGZmM2NjIiwidHlwIjoiQmVhcmVyIiwiYXpwIjoiSXN0aW8iLCJzZXNzaW9uX3N0YXRlIjoiMzkyMWU0NWMtOWNkNS00ZWE5LWFjZTktMTZhM2MyMzUzMzdiIiwiYWNyIjoiMSIsImFsbG93ZWQtb3JpZ2lucyI6WyIvKiJdLCJyZWFsbV9hY2Nlc3MiOnsicm9sZXMiOlsib2ZmbGluZV9hY2Nlc3MiLCJhZG1pbiIsInVtYV9hdXRob3JpemF0aW9uIiwiZGVmYXVsdC1yb2xlcy1pc3RpbyJdfSwicmVzb3VyY2VfYWNjZXNzIjp7ImFjY291bnQiOnsicm9sZXMiOlsibWFuYWdlLWFjY291bnQiLCJtYW5hZ2UtYWNjb3VudC1saW5rcyIsInZpZXctcHJvZmlsZSJdfX0sInNjb3BlIjoicHJvZmlsZSBlbWFpbCIsInNpZCI6IjM5MjFlNDVjLTljZDUtNGVhOS1hY2U5LTE2YTNjMjM1MzM3YiIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJuYW1lIjoiTWlsYW4gQWpkZXIiLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJob3RlbC1hZG1pbiIsImdpdmVuX25hbWUiOiJNaWxhbiIsImZhbWlseV9uYW1lIjoiQWpkZXIiLCJlbWFpbCI6ImFqZGVyLm1pbGFuMjAwMEBnbWFpbC5jb20ifQ.JBPiqSMdgcuuhDu8tJ3RAm7Vq0TVZD0KrcuKvObRXKuC7XsKAfQVX0mvd7JZYsmjgIdOMHYC6f8HxC-sXbdXvgzAJeCyx9icW0hfYQ1qfuizW_AfFnTbmqovZFeJME6c2n_okUKM8bNXOHWNvVR0uZ6UdrBiyNemNfXRjSKA8SAMiIOB_ZbH4aUqEtrNbZ3CEe3G8ZMUJ8UlYdzUeBt2VxPZenBqwCkhmYAx4pmdkuF5f0pB8sMKsbpjgkssgYnLwalg9tDRmxLNOv7eeBenyiZm73EwFfnE0xvf0KtmauwYBUkx91AFZfuvlRBM1u4O8PGfo24iMV_S524GfbSMsA"
curl -X PUT http://hotel/hotel/accommodation/price/66467094ea31a50a941a503d -d '{"price": 555.50, "type": "PerApartmentUnit"}' -H "Content-Type: application/json" -H "Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJSM19wNUs5dFVtTkpmSVlhdmo1NV8xTjI1VUJOMmh5cmJJWGVnVmV4ZG1ZIn0.eyJleHAiOjE3MTYwMzM2MTgsImlhdCI6MTcxNjAzMzMxOCwianRpIjoiNGFkMmZmM2ItYTAwYS00NmQxLWE4ZDktNDNjNjJmNWUyMTA3IiwiaXNzIjoiaHR0cDovL2tleWNsb2FrLmRlZmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWw6ODA4MC9yZWFsbXMvSXN0aW8iLCJhdWQiOiJhY2NvdW50Iiwic3ViIjoiNTY5YjkxMjktN2NjYi00OWNhLTgzY2YtMDZmNmMzMGZmM2NjIiwidHlwIjoiQmVhcmVyIiwiYXpwIjoiSXN0aW8iLCJzZXNzaW9uX3N0YXRlIjoiMzkyMWU0NWMtOWNkNS00ZWE5LWFjZTktMTZhM2MyMzUzMzdiIiwiYWNyIjoiMSIsImFsbG93ZWQtb3JpZ2lucyI6WyIvKiJdLCJyZWFsbV9hY2Nlc3MiOnsicm9sZXMiOlsib2ZmbGluZV9hY2Nlc3MiLCJhZG1pbiIsInVtYV9hdXRob3JpemF0aW9uIiwiZGVmYXVsdC1yb2xlcy1pc3RpbyJdfSwicmVzb3VyY2VfYWNjZXNzIjp7ImFjY291bnQiOnsicm9sZXMiOlsibWFuYWdlLWFjY291bnQiLCJtYW5hZ2UtYWNjb3VudC1saW5rcyIsInZpZXctcHJvZmlsZSJdfX0sInNjb3BlIjoicHJvZmlsZSBlbWFpbCIsInNpZCI6IjM5MjFlNDVjLTljZDUtNGVhOS1hY2U5LTE2YTNjMjM1MzM3YiIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJuYW1lIjoiTWlsYW4gQWpkZXIiLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJob3RlbC1hZG1pbiIsImdpdmVuX25hbWUiOiJNaWxhbiIsImZhbWlseV9uYW1lIjoiQWpkZXIiLCJlbWFpbCI6ImFqZGVyLm1pbGFuMjAwMEBnbWFpbC5jb20ifQ.JBPiqSMdgcuuhDu8tJ3RAm7Vq0TVZD0KrcuKvObRXKuC7XsKAfQVX0mvd7JZYsmjgIdOMHYC6f8HxC-sXbdXvgzAJeCyx9icW0hfYQ1qfuizW_AfFnTbmqovZFeJME6c2n_okUKM8bNXOHWNvVR0uZ6UdrBiyNemNfXRjSKA8SAMiIOB_ZbH4aUqEtrNbZ3CEe3G8ZMUJ8UlYdzUeBt2VxPZenBqwCkhmYAx4pmdkuF5f0pB8sMKsbpjgkssgYnLwalg9tDRmxLNOv7eeBenyiZm73EwFfnE0xvf0KtmauwYBUkx91AFZfuvlRBM1u4O8PGfo24iMV_S524GfbSMsA"
curl -X PUT http://hotel/hotel/accommodation/price/66467094ea31a50a941a503d -d '{"price": 555.50, "type": "PerApartmentUnit"}' -H "Content-Type: application/json" -H "Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJSM19wNUs5dFVtTkpmSVlhdmo1NV8xTjI1VUJOMmh5cmJJWGVnVmV4ZG1ZIn0.eyJleHAiOjE3MTYwMzM2MTgsImlhdCI6MTcxNjAzMzMxOCwianRpIjoiNGFkMmZmM2ItYTAwYS00NmQxLWE4ZDktNDNjNjJmNWUyMTA3IiwiaXNzIjoiaHR0cDovL2tleWNsb2FrLmRlZmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWw6ODA4MC9yZWFsbXMvSXN0aW8iLCJhdWQiOiJhY2NvdW50Iiwic3ViIjoiNTY5YjkxMjktN2NjYi00OWNhLTgzY2YtMDZmNmMzMGZmM2NjIiwidHlwIjoiQmVhcmVyIiwiYXpwIjoiSXN0aW8iLCJzZXNzaW9uX3N0YXRlIjoiMzkyMWU0NWMtOWNkNS00ZWE5LWFjZTktMTZhM2MyMzUzMzdiIiwiYWNyIjoiMSIsImFsbG93ZWQtb3JpZ2lucyI6WyIvKiJdLCJyZWFsbV9hY2Nlc3MiOnsicm9sZXMiOlsib2ZmbGluZV9hY2Nlc3MiLCJhZG1pbiIsInVtYV9hdXRob3JpemF0aW9uIiwiZGVmYXVsdC1yb2xlcy1pc3RpbyJdfSwicmVzb3VyY2VfYWNjZXNzIjp7ImFjY291bnQiOnsicm9sZXMiOlsibWFuYWdlLWFjY291bnQiLCJtYW5hZ2UtYWNjb3VudC1saW5rcyIsInZpZXctcHJvZmlsZSJdfX0sInNjb3BlIjoicHJvZmlsZSBlbWFpbCIsInNpZCI6IjM5MjFlNDVjLTljZDUtNGVhOS1hY2U5LTE2YTNjMjM1MzM3YiIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJuYW1lIjoiTWlsYW4gQWpkZXIiLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJob3RlbC1hZG1pbiIsImdpdmVuX25hbWUiOiJNaWxhbiIsImZhbWlseV9uYW1lIjoiQWpkZXIiLCJlbWFpbCI6ImFqZGVyLm1pbGFuMjAwMEBnbWFpbC5jb20ifQ.JBPiqSMdgcuuhDu8tJ3RAm7Vq0TVZD0KrcuKvObRXKuC7XsKAfQVX0mvd7JZYsmjgIdOMHYC6f8HxC-sXbdXvgzAJeCyx9icW0hfYQ1qfuizW_AfFnTbmqovZFeJME6c2n_okUKM8bNXOHWNvVR0uZ6UdrBiyNemNfXRjSKA8SAMiIOB_ZbH4aUqEtrNbZ3CEe3G8ZMUJ8UlYdzUeBt2VxPZenBqwCkhmYAx4pmdkuF5f0pB8sMKsbpjgkssgYnLwalg9tDRmxLNOv7eeBenyiZm73EwFfnE0xvf0KtmauwYBUkx91AFZfuvlRBM1u4O8PGfo24iMV_S524GfbSMsA"
curl -X PUT http://hotel/hotel/accommodation/price/66467094ea31a50a941a503e -d '{"price": 12555.50}' -H "Content-Type: application/json" -H "Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJSM19wNUs5dFVtTkpmSVlhdmo1NV8xTjI1VUJOMmh5cmJJWGVnVmV4ZG1ZIn0.eyJleHAiOjE3MTYwMzM2MTgsImlhdCI6MTcxNjAzMzMxOCwianRpIjoiNGFkMmZmM2ItYTAwYS00NmQxLWE4ZDktNDNjNjJmNWUyMTA3IiwiaXNzIjoiaHR0cDovL2tleWNsb2FrLmRlZmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWw6ODA4MC9yZWFsbXMvSXN0aW8iLCJhdWQiOiJhY2NvdW50Iiwic3ViIjoiNTY5YjkxMjktN2NjYi00OWNhLTgzY2YtMDZmNmMzMGZmM2NjIiwidHlwIjoiQmVhcmVyIiwiYXpwIjoiSXN0aW8iLCJzZXNzaW9uX3N0YXRlIjoiMzkyMWU0NWMtOWNkNS00ZWE5LWFjZTktMTZhM2MyMzUzMzdiIiwiYWNyIjoiMSIsImFsbG93ZWQtb3JpZ2lucyI6WyIvKiJdLCJyZWFsbV9hY2Nlc3MiOnsicm9sZXMiOlsib2ZmbGluZV9hY2Nlc3MiLCJhZG1pbiIsInVtYV9hdXRob3JpemF0aW9uIiwiZGVmYXVsdC1yb2xlcy1pc3RpbyJdfSwicmVzb3VyY2VfYWNjZXNzIjp7ImFjY291bnQiOnsicm9sZXMiOlsibWFuYWdlLWFjY291bnQiLCJtYW5hZ2UtYWNjb3VudC1saW5rcyIsInZpZXctcHJvZmlsZSJdfX0sInNjb3BlIjoicHJvZmlsZSBlbWFpbCIsInNpZCI6IjM5MjFlNDVjLTljZDUtNGVhOS1hY2U5LTE2YTNjMjM1MzM3YiIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJuYW1lIjoiTWlsYW4gQWpkZXIiLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJob3RlbC1hZG1pbiIsImdpdmVuX25hbWUiOiJNaWxhbiIsImZhbWlseV9uYW1lIjoiQWpkZXIiLCJlbWFpbCI6ImFqZGVyLm1pbGFuMjAwMEBnbWFpbC5jb20ifQ.JBPiqSMdgcuuhDu8tJ3RAm7Vq0TVZD0KrcuKvObRXKuC7XsKAfQVX0mvd7JZYsmjgIdOMHYC6f8HxC-sXbdXvgzAJeCyx9icW0hfYQ1qfuizW_AfFnTbmqovZFeJME6c2n_okUKM8bNXOHWNvVR0uZ6UdrBiyNemNfXRjSKA8SAMiIOB_ZbH4aUqEtrNbZ3CEe3G8ZMUJ8UlYdzUeBt2VxPZenBqwCkhmYAx4pmdkuF5f0pB8sMKsbpjgkssgYnLwalg9tDRmxLNOv7eeBenyiZm73EwFfnE0xvf0KtmauwYBUkx91AFZfuvlRBM1u4O8PGfo24iMV_S524GfbSMsA"
curl -X PUT http://hotel/hotel/accommodation/price/66467094ea31a50a941a503e -d '{}' -H "Content-Type: application/json" -H "Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJSM19wNUs5dFVtTkpmSVlhdmo1NV8xTjI1VUJOMmh5cmJJWGVnVmV4ZG1ZIn0.eyJleHAiOjE3MTYwMzM2MTgsImlhdCI6MTcxNjAzMzMxOCwianRpIjoiNGFkMmZmM2ItYTAwYS00NmQxLWE4ZDktNDNjNjJmNWUyMTA3IiwiaXNzIjoiaHR0cDovL2tleWNsb2FrLmRlZmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWw6ODA4MC9yZWFsbXMvSXN0aW8iLCJhdWQiOiJhY2NvdW50Iiwic3ViIjoiNTY5YjkxMjktN2NjYi00OWNhLTgzY2YtMDZmNmMzMGZmM2NjIiwidHlwIjoiQmVhcmVyIiwiYXpwIjoiSXN0aW8iLCJzZXNzaW9uX3N0YXRlIjoiMzkyMWU0NWMtOWNkNS00ZWE5LWFjZTktMTZhM2MyMzUzMzdiIiwiYWNyIjoiMSIsImFsbG93ZWQtb3JpZ2lucyI6WyIvKiJdLCJyZWFsbV9hY2Nlc3MiOnsicm9sZXMiOlsib2ZmbGluZV9hY2Nlc3MiLCJhZG1pbiIsInVtYV9hdXRob3JpemF0aW9uIiwiZGVmYXVsdC1yb2xlcy1pc3RpbyJdfSwicmVzb3VyY2VfYWNjZXNzIjp7ImFjY291bnQiOnsicm9sZXMiOlsibWFuYWdlLWFjY291bnQiLCJtYW5hZ2UtYWNjb3VudC1saW5rcyIsInZpZXctcHJvZmlsZSJdfX0sInNjb3BlIjoicHJvZmlsZSBlbWFpbCIsInNpZCI6IjM5MjFlNDVjLTljZDUtNGVhOS1hY2U5LTE2YTNjMjM1MzM3YiIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJuYW1lIjoiTWlsYW4gQWpkZXIiLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJob3RlbC1hZG1pbiIsImdpdmVuX25hbWUiOiJNaWxhbiIsImZhbWlseV9uYW1lIjoiQWpkZXIiLCJlbWFpbCI6ImFqZGVyLm1pbGFuMjAwMEBnbWFpbC5jb20ifQ.JBPiqSMdgcuuhDu8tJ3RAm7Vq0TVZD0KrcuKvObRXKuC7XsKAfQVX0mvd7JZYsmjgIdOMHYC6f8HxC-sXbdXvgzAJeCyx9icW0hfYQ1qfuizW_AfFnTbmqovZFeJME6c2n_okUKM8bNXOHWNvVR0uZ6UdrBiyNemNfXRjSKA8SAMiIOB_ZbH4aUqEtrNbZ3CEe3G8ZMUJ8UlYdzUeBt2VxPZenBqwCkhmYAx4pmdkuF5f0pB8sMKsbpjgkssgYnLwalg9tDRmxLNOv7eeBenyiZm73EwFfnE0xvf0KtmauwYBUkx91AFZfuvlRBM1u4O8PGfo24iMV_S524GfbSMsA"


eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJSM19wNUs5dFVtTkpmSVlhdmo1NV8xTjI1VUJOMmh5cmJJWGVnVmV4ZG1ZIn0.eyJleHAiOjE3MTYwMzM2MTgsImlhdCI6MTcxNjAzMzMxOCwianRpIjoiNGFkMmZmM2ItYTAwYS00NmQxLWE4ZDktNDNjNjJmNWUyMTA3IiwiaXNzIjoiaHR0cDovL2tleWNsb2FrLmRlZmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWw6ODA4MC9yZWFsbXMvSXN0aW8iLCJhdWQiOiJhY2NvdW50Iiwic3ViIjoiNTY5YjkxMjktN2NjYi00OWNhLTgzY2YtMDZmNmMzMGZmM2NjIiwidHlwIjoiQmVhcmVyIiwiYXpwIjoiSXN0aW8iLCJzZXNzaW9uX3N0YXRlIjoiMzkyMWU0NWMtOWNkNS00ZWE5LWFjZTktMTZhM2MyMzUzMzdiIiwiYWNyIjoiMSIsImFsbG93ZWQtb3JpZ2lucyI6WyIvKiJdLCJyZWFsbV9hY2Nlc3MiOnsicm9sZXMiOlsib2ZmbGluZV9hY2Nlc3MiLCJhZG1pbiIsInVtYV9hdXRob3JpemF0aW9uIiwiZGVmYXVsdC1yb2xlcy1pc3RpbyJdfSwicmVzb3VyY2VfYWNjZXNzIjp7ImFjY291bnQiOnsicm9sZXMiOlsibWFuYWdlLWFjY291bnQiLCJtYW5hZ2UtYWNjb3VudC1saW5rcyIsInZpZXctcHJvZmlsZSJdfX0sInNjb3BlIjoicHJvZmlsZSBlbWFpbCIsInNpZCI6IjM5MjFlNDVjLTljZDUtNGVhOS1hY2U5LTE2YTNjMjM1MzM3YiIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJuYW1lIjoiTWlsYW4gQWpkZXIiLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJob3RlbC1hZG1pbiIsImdpdmVuX25hbWUiOiJNaWxhbiIsImZhbWlseV9uYW1lIjoiQWpkZXIiLCJlbWFpbCI6ImFqZGVyLm1pbGFuMjAwMEBnbWFpbC5jb20ifQ.JBPiqSMdgcuuhDu8tJ3RAm7Vq0TVZD0KrcuKvObRXKuC7XsKAfQVX0mvd7JZYsmjgIdOMHYC6f8HxC-sXbdXvgzAJeCyx9icW0hfYQ1qfuizW_AfFnTbmqovZFeJME6c2n_okUKM8bNXOHWNvVR0uZ6UdrBiyNemNfXRjSKA8SAMiIOB_ZbH4aUqEtrNbZ3CEe3G8ZMUJ8UlYdzUeBt2VxPZenBqwCkhmYAx4pmdkuF5f0pB8sMKsbpjgkssgYnLwalg9tDRmxLNOv7eeBenyiZm73EwFfnE0xvf0KtmauwYBUkx91AFZfuvlRBM1u4O8PGfo24iMV_S524GfbSMsA
```
Get JWT token for user
```shell
curl -X POST -d "client_id=Istio" -d "username=hotel-user" -d "password=test" -d "grant_type=password" "http://keycloak.default.svc.cluster.local:8080/realms/Istio/protocol/openid-connect/token"
```

Get JWT token for admin
```shell
 curl -X POST -d "client_id=Istio" -d "username=hotel-admin" -d "password=test" -d "grant_type=password" "http://keycloak.default.svc.cluster.local:8080/realms/Istio/protocol/openid-connect/token"
```

Ingress setup:
Deploy ingress:
```shell
kubectl -n backend apply -f ingress.yml
kubectl -n backend describe ingress demo-ingress
```

Apply za ceo ili vise direktorijuma
```shell
kubectl -n backend  apply -R -f k8s
kubectl -n backend  apply -R -f istio
```

Ponisti prethodnu verziju i apply novu 
```shell
kubectl replace --force -f ingress.yml
kubectl replace --force -f istio/authorizationPolicy.yaml
kubectl replace --force -f k8s/booking-service.yml
kubectl -n backend  replace --force -f k8s/booking-configmap.yml
```

Keycloak
```shell
minikube addons enable ingress
kubectl create -f https://raw.githubusercontent.com/keycloak/keycloak-quickstarts/latest/kubernetes/keycloak.yaml
minikube tunnel

browser: localhost:8080 (username: admin, password: admin)

Create Istio realm
Create Istio client 
Create hotel-user , hotel-admin (password: test)
```


docker build -t devopszms2024/zms-devops-booking-service .
docker push devopszms2024/zms-devops-booking-service
kubectl replace --force -f k8s/
# regAccess
docker registry access control


## 运行 

```
./regAccess -port 10080 -htpasswd auth/htpasswd
```

## 新建 auth user/passwd for docker registry：

```
$ mkdir auth

$ docker run \
  --entrypoint htpasswd \
  registry:2 -Bbn testuser testpassword >> auth/htpasswd


$ docker run --rm \
  --entrypoint htpasswd \
  registry:2 -Bbn holly 123456 >> auth/htpasswd
```

## 运行 docker registry

```
 $ cd /path/to/regAccess/

 $ docker run -d \
  -p 5000:5000 \
  --restart=always \
  --name registry \
  -v `pwd`/auth:/auth \
  -e "REGISTRY_AUTH=htpasswd" \
  -e "REGISTRY_AUTH_HTPASSWD_REALM=Registry Realm" \
  -e REGISTRY_AUTH_HTPASSWD_PATH=/auth/htpasswd \
  -v `pwd`/certs:/certs \
  -e REGISTRY_HTTP_TLS_CERTIFICATE=/certs/myreg.crt \
  -e REGISTRY_HTTP_TLS_KEY=/certs/myreg.key \
  registry:2
  
```
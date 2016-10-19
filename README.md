#Abattoir
Cleans up dead rancher hosts from rancher provided they no longer exist, currently depends on a custom API that provides a list of servers.

##Configuration

Requires the following environment variables to run
```
export ABATTOIR_RANCHER_ACCESS_KEY=<rancher access key>
export ABATTOIR_RANCHER_SECRET_KEY=<rancher secret key>
export ABATTOIR_RANCHER_URL=https://rancher.somedomain.com
export ABATTOIR_SERVER_API=https://ops-api-servers.somehwere.com
export ABATTOIR_RANCHER_PROJECT_ID=1a5
```

##Running

```
docker-compose up
```

##Build

```
make all
```

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
export ABATTOIR_RUN_INTERVAL=300
```

ABATTOIR_RANCHER_ACCESS_KEY
* Access key for rancher

ABATTOIR_RANCHER_SECRET_KEY
* Secret key for rancher

ABATTOIR_RANCHER_URL
* Can be in the format http://rancher.somewhere.com:8080 or https://rancher.somewhere.com

ABATTOIR_RANCHER_PROJECT_ID
* This is your rancher project id, it's 1a5 for the default env, but might be different depending on your setup

ABATTOIR_SERVER_API
* The API that abattoir connects to in order to get a list of running hosts for your current environment.  It's used to verify that the host which is showing down in rancher is gone from your infrastructure

ABATTOIR_RUN_INTERVAL
* How often abattoir should poll rancher for inactive hosts


##Running

```
docker-compose up
```

##Build

```
make all
```

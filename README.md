<p align="center">
  <a href="https://nuntio.io">
    <picture>
      <img alt="Nuntio logo" src="https://github.com/nuntiodev/nuntio-go-api/assets/22043/d5981401-a9ac-42e4-b584-1b2718a47ea9" width="40%">
    </picture>
  </a>
</p>
<h2 align="center">
  Nuntio [feature] demo
</h2>

<h4 align="center">
  <a href="...">Documentation</a>
</h4>
<p align="center">
  <a href="[https://twitter.com/intent/follow?screen_name=nuntio_io](https://discord.com/invite/9KYQqpRpBN)">
    <img src="https://img.shields.io/discord/1076063204893012049" alt="Join us on Discord" />
  </a>
  <a href="https://twitter.com/intent/follow?screen_name=nuntio_io">
    <img src="https://img.shields.io/twitter/follow/nuntio_io?label=Follow%20@nuntio_io" alt="Follow @nuntio_io" />
  </a>
</p>

## Nuntio [feature] demo

This demo showcases how Nuntio helps with managing databases. Specifically it showcases Nuntio operating with a MongoDB database.

## What is Nuntio?

Nuntio is an open-source cloud development platform for Kubernetes. It features simple-to-use Capsule for Application deployments and batteries-included Modules for Auth, User-management, Storage and Databases.

## Prerequisites

Nuntio must be running either locally in Docker or on a Kubernetes cluster. You can refer to [this guide](https://beta-docs.nuntio.io/get-started) to learn how to install Nuntio.

## Get started

To run the demo locally, follow these steps:

### Step 1. Create a new Nuntio managed MongoDB database

```
nuntio database create --name our_db --type mongo
```

### Step 2: Add credentials to the newly created database

```
nuntio database create-credentials
```

And write `our_db` once it prompts for the DB identifier.
Store the `clientID` and `secret` it outputs somewhere you can access later.

### Step 3: Clone this repo and build a Docker image from it

```
git clone git@github.com:nuntiodev/database-demo.git
cd database-demo
docker build -t database-demo .
```

### Step 4: Create a new capsule and deploy the database-demo docker image to it

```
nuntio capsule create --name database-demo
nuntio capsule create-build database-demo --image --database-demo --deploy
```

### Step 5: Automatically inject Nuntio client credentials and database credentials

```
nuntio capsule config database-demo --add-credentials
```

From the dashboard under the `settings` tab for your capsule, set the following environment variables from which the demo expects to read database name and credentials
![image](https://i.imgur.com/LAIaB1E.png)

### Step 6: Expose the port to the public which the application listens to

`Network` in your capsule's dashboard
![image](https://i.imgur.com/lAHNeA7.png)

You should now be able to access the endpoints which the capsule implements on `localhost:3333/`

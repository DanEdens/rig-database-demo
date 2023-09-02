<p align="center">
  <a href="https://rig.dev">
    <picture>
      <img alt="Rig logo" src="https://github-production-user-asset-6210df.s3.amazonaws.com/22043/265044731-d8433e7b-4fe5-40ee-a2a2-bbfec9b26c63.svg" width="40%">
    </picture>
  </a>
</p>
<h2 align="center">
  Rig database demo
</h2>

<h4 align="center">
  <a href="https://docs.rig.dev">Documentation</a>
</h4>
<p align="center">
  <a href="[https://twitter.com/intent/follow?screen_name=rig_dev](https://discord.gg/Tn5wmXMM2U)">
    <img src="https://img.shields.io/discord/1076063204893012049" alt="Join us on Discord" />
  </a>
  <a href="https://twitter.com/intent/follow?screen_name=rig_dev">
    <img src="https://img.shields.io/twitter/follow/rig_dev?label=Follow%20@rig_dev" alt="Follow @rig_dev" />
  </a>
</p>

## Rig database demo

This demo showcases how Rig helps with managing databases. Specifically it showcases Rig operating with a MongoDB database. See [here](http://docs.rig.dev) for a more comprehensive walkthrough of the demo.

## What is Rig?

Rig is an open-source cloud development platform for Kubernetes. It features simple-to-use Capsule for Application deployments and batteries-included Modules for Auth, User-management, Storage and Databases.

## Prerequisites

Rig must be running either locally in Docker or on a Kubernetes cluster. You can refer to [this guide](https://docs.rig.dev/get-started) to learn how to install Rig.

## Get started

To run the demo locally, follow these steps:

### Step 1. Create a new Rig managed MongoDB database

```
rig database create --name our_db --type mongo
```

### Step 2: Add credentials to the newly created database

```
rig database create-credentials
```

And write `our_db` once it prompts for the DB identifier.
Store the `clientID` and `secret` it outputs somewhere you can access later.

### Step 3: Clone this repo and build a Docker image from it

```
git clone git@github.com:rigdev/database-demo.git
cd database-demo
docker build -t database-demo .
```

### Step 4: Create a new capsule and deploy the database-demo docker image to it

```
rig capsule create --name database-demo
rig capsule create-build database-demo --image --database-demo --deploy
```

### Step 5: Automatically inject Rig client credentials and database credentials

```
rig capsule config database-demo --add-credentials
```

Also add credentials to the database we created and copy the secret

```
rig database create-credentials our_db
```

From the dashboard under the `settings` tab for your capsule, set the following environment variables from which the demo expects to read database name and credentials
![image](https://i.imgur.com/LAIaB1E.png)

### Step 6: Expose the port to the public which the application listens to

`Network` in your capsule's dashboard
![image](https://i.imgur.com/lAHNeA7.png)

You should now be able to access the endpoints which the capsule implements on `localhost:3333/`

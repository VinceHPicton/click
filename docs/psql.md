# Using Postgres

This section describe operations around the use of Postgres in relation to our project.

## Runing Postgres locally

Use the script `./scripts/psql.sh` to work with local setup.

This is typical sequence of script use cases:

* STEP 1 - To build or pull require image run the command `./scripts/psql.sh image build`. This will pull a PSQL server image, a pgadmin image and build a psql cli image.

* STEP 2 - Start a docker-compose network, run the command `./scripts/psql.sh network start`. The script for this is found in [./deployment/compose/docker-compose.yml](../deployment/compose/docker-compose.yml)

* STEP 3 - Assuming you have a running network, you will find a temporary folder `./tmp/pgdata` where your psql data is found. You can also open PGAdmin console via your browser `http://localhost:5050` (assuming you didn't change the port number).

* STEP 4 - To interact with the psql server there is a default cli. Run the command `./script/pssql.sh client cli`. In the shell, you can run the command `psql` which will open up a shell based DB scripting solution.

* STEP 5 - To shut the network, run the command `./scripts/psql.sh network stop`.

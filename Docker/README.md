
# Integration rextporter and skycoin using Docker

## Steps guide
### 1- Build `rextporter` image
```
docker build -t rextporter ./Docker/Dockerfile
```
### 2- Create volume `rextporter-config`
```
docker volume create rextporter-config
```
This volume will contain the rextporter config files, if you want to specified a directory containing the config files you can use yhe next command instead of the previous:
```
docker volume create --driver local \
                     --opt type=none \
                     --opt o=bind \
                     --opt device=/path/to/config-directory \
                     rextporter-config
```
### 3- Run docker-compose
```
cd ./Docker/docker-compose.yml
docker-compose up
```

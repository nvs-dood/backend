### Setting up mariadb docker container for use with backend
docker run -p 3306:3306 --name dood-mariadb -e MYSQL_ROOT_PASSWORD=Pass -d mariadb:latest

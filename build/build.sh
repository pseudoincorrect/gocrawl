#! /bin/bash

ROOT_PATH="$PWD"/..
GO_SERVICES=(worker server)
OTHER_SERVICES=(rabbitmq)
ALL_SERVICES=("${GO_SERVICES[@]}" "${OTHER_SERVICES[@]}")

help() {
  echo "
Usage: ./build COMMAND [OPTIONS]

Tool to help building and deploying GoCrawl

Availble commands:

  build-all                 Build all services

  build-one [service_name]  Build one service

  deploy-all                Deploy the application (docker-compose)

  down-all                  Shut the down the application

  down-one [service_name]   Shut one service

  sbd-one [service_name]    Shutdown, build and deploy one service

  log-one [service_name]    Get the logs of a service

  log-all                   Get the logs for all services and follow

  help:                     Print this page
"
}

check_go_service() {
  if [[ ! " ${GO_SERVICES[@]} " =~ " $1 " ]]; then
    echo "Service (Go) $1 does not exist, aborting..."
    exit
  fi
}

check_all_service() {
  if [[ ! " ${ALL_SERVICES[@]} " =~ " $1 " ]]; then
    echo "Service $1 does not exist, aborting..."
    exit
  fi
}

build_all() {
  echo "Building all services"
  for service in "${GO_SERVICES[@]}"; do
    build_one $service
  done
}

build_one() {
  service=$1
  check_go_service $service
  echo "Building service $service"
  cd "$ROOT_PATH"
  docker build . -t $service -f ./docker/dockerfile --build-arg SVC=$service
}

deploy_all() {
  echo "Deploying the application"
  cd "$ROOT_PATH/docker"
  docker-compose up -d
}

dc_down_all() {
  echo "Shutting down the application"
  cd "$ROOT_PATH/docker"
  docker-compose down
}

dc_down_one() {
  service=$1
  check_all_service $service
  echo "Shutting down service $service"
  cd "$ROOT_PATH/docker"
  docker-compose stop $service
  docker-compose rm --force $service
}

shutdown_build_deploy_one() {
  service=$1
  check_go_service $service
  echo "shutdowm, build, deploy service $service"
  cd "$ROOT_PATH/docker"
  docker-compose stop $service
  docker-compose rm --force $service
  build_one $service
  cd "$ROOT_PATH/docker"
  docker-compose up -d
}

log_one() {
  service=$1
  check_go_service $service
  echo "getting logs for service $service"
  cd "$ROOT_PATH/docker"
  echo \"$service\"
  docker-compose logs | grep -E "$service"
}

log_all() {
  echo "getting all logs and following log output"
  cd "$ROOT_PATH/docker"
  docker-compose logs -f
}

##########

case $1 in
help)
  help
  ;;
build-all)
  build_all
  ;;
build-one)
  build_one $2
  ;;
deploy-all)
  deploy_all
  ;;
down-all)
  dc_down_all
  ;;
down-one)
  dc_down_one $2
  ;;
sbd-one)
  shutdown_build_deploy_one $2
  ;;
log-one)
  log_one $2
  ;;
log-all)
  log_all $2
  ;;
*)
  echo "unknown command \"$1\""
  help
  ;;
esac

#!?usr/bin/env bash
cd ansible
ansible-playbook -i "localhost," -c local site.yml
cd ../getweather
docker build -t weather:dev .
docker run --rm -e OPENWEATHER_API_KEY="xxxxxxxxxxxx" -e CITY_NAME="Honolulu" weather:dev

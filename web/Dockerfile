FROM nginx

RUN adduser --disabled-password --gecos '' web
USER web

COPY . /usr/share/nginx/html

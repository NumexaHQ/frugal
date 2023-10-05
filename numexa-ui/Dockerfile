# stage1 as builder
FROM node:lts as builder

# Set working directory
WORKDIR /react-ui
# Copy all files from current directory to working dir in image
COPY . .
# install node modules and build assets
RUN yarn install && yarn build


FROM nginx:alpine

#!/bin/sh

# Set working directory to nginx asset directory
WORKDIR /usr/share/nginx/html

## Remove default nginx index page
RUN rm -rf ./*

# Copy from the stage 1
COPY .nginx/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /react-ui/build .

EXPOSE 3000 80

ENTRYPOINT ["nginx", "-g", "daemon off;"]

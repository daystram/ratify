FROM node:lts-alpine as builder
WORKDIR /build
#TODO: add .env.production when ready
COPY package.json yarn.lock ./
RUN yarn install
COPY . .
RUN yarn build --mode production

FROM nginx:stable-alpine as production-stage
COPY --from=builder /build/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]

FROM node:11.3 AS builder

WORKDIR /usr/src/app

COPY backend/package*.json backend/
COPY backend/src backend/src
COPY backend/tsconfig.json backend/

RUN cd backend && npm install
RUN cd backend && npm run tsc:build

COPY frontend/package*.json frontend/
COPY frontend/src frontend/src
COPY frontend/tsconfig*.json frontend/
COPY frontend/angular.json frontend/
COPY frontend/karma.conf.js frontend/

RUN cd frontend && npm install
RUN cd frontend && npm run build

RUN mkdir app && cp -rf backend/dist app && ls app
RUN cp -rf frontend/dist/uploader app/public && ls app && ls app/public

###

FROM node:11.3 AS installer

WORKDIR /usr/src/app

COPY --from=builder /usr/src/app/app/dist .
COPY --from=builder /usr/src/app/app/public ./public
COPY --from=builder /usr/src/app/backend/package*.json ./

RUN npm install --production
# ---

FROM node:11.3-slim

WORKDIR /app

COPY --from=installer /usr/src/app ./
COPY --from=installer /usr/src/app/node_modules node_modules

ENTRYPOINT ["node"]
CMD ["index.js"]
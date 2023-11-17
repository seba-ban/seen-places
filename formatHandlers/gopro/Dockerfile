FROM node:20.9.0-slim as builder

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY . .

RUN npm run build

FROM builder AS dev

ENV NODE_ENV=development

CMD [ "npm", "run", "dev" ]

FROM node:20.9.0-slim as production

ENV NODE_ENV=production
WORKDIR /app

RUN apt-get update && apt-get install ffmpeg -y

COPY package*.json ./
RUN npm install --omit=dev

COPY --from=builder /app/dist ./dist

VOLUME [ "/data" ]
ENV LOCAL_STORAGE_FILES_DIR /data/storage

CMD ["npm", "run", "start"]

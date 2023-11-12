FROM node:20.9.0-slim as builder

WORKDIR /builder

COPY package*.json ./

RUN npm install

COPY . .

RUN npm run build

FROM node:20.9.0-slim as production

ENV NODE_ENV=production
WORKDIR /app

RUN apt-get update && apt-get install ffmpeg -y

COPY package*.json ./
RUN npm install --omit=dev

COPY --from=builder /builder/dist ./dist

CMD ["npm", "run", "start"]

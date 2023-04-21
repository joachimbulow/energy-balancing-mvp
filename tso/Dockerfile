FROM node:19-alpine

WORKDIR /app

COPY ["package.json", "package-lock.json*", "./"]

ENV REDIS_HOST=redis

RUN npm install --production

COPY . .

CMD ["node", "main.js"]
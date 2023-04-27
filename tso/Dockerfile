FROM node:19-alpine

WORKDIR /app

COPY ["package.json", "package-lock.json*", "./"]

ENV BROKER=kafka
ENV BROKER_URL=20.229.248.72
ENV BROKER_PORT=9092

RUN npm install --production

COPY . .

CMD ["node", "main.js"]
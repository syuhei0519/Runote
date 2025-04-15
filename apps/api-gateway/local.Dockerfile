FROM node:22

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

ENV NODE_ENV=develop

CMD ["npx", "ts-node-dev", "--respawn", "src/index.ts"]
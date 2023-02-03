FROM node:14-bullseye

RUN mkdir /app
# Create app directory
WORKDIR /app

# Install app dependencies
# A wildcard is used to ensure both package.json AND package-lock.json are copied
# where available (npm@5+)
COPY package.json /app

RUN npm install
COPY . /app
RUN cd /app/client && npm install && npm run build && cd /app;
RUN cd /app/server && npm install && cd /app;
# If you are building your code for production
# RUN npm ci --only=production
RUN ls
# Bundle app source


EXPOSE 3000
CMD [ "npm", "run", "s2" ]

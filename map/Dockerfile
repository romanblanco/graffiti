FROM node:16.0.0-alpine
WORKDIR /app
ENV PATH /app/node_modules/.bin:$PATH
COPY . ./
# reason: getaddrinfo EAI_AGAIN registry.npmjs.org
RUN npm config set registry http://registry.yarnpkg.com
RUN npm --verbose install
RUN npm --verbose install -g react-scripts@latest
CMD ["npm", "start"]

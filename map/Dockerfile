FROM node:16.0.0-alpine
WORKDIR /app
ENV PATH /app/node_modules/.bin:$PATH
COPY . ./
RUN npm --verbose install
RUN npm install -g react-scripts@latest
CMD ["npm", "start"]

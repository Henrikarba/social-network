# Use the long-term support (LTS) version of the Node.js Alpine image for the build stage
FROM node:lts-alpine AS build

# Set the working directory inside the container
WORKDIR /app

# Copy package.json and package-lock.json to the container's working directory
COPY package.json ./
COPY package-lock.json ./

# Install the Node.js dependencies inside the container
RUN npm install

# Copy the rest of the application code into the container
COPY . ./

# Build the application inside the container
RUN npm run build

# Start a new stage using the lightweight busybox:1.35 image
FROM busybox:1.35

# Create a non-root user to own the files and run our server
RUN adduser -D static
USER static
WORKDIR /home/static

# Copy the built application from the previous stage into the current directory
COPY --from=build /app/dist .

# Run BusyBox httpd
CMD echo "Serving frontend on: localhost:5000" && busybox httpd -f -v -p 5000

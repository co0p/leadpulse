# Stage 1: Build
FROM node:20-alpine AS build

WORKDIR /build

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy source
COPY . .

# Build Vue app
RUN npm run build

# Stage 2: Runtime (serve static files)
FROM nginx:alpine

# Copy built app from build stage
COPY --from=build /build/dist /usr/share/nginx/html

# Copy Nginx config to proxy API calls
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose port 3000
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
	CMD wget --no-verbose --tries=1 --spider http://localhost:3000/ || exit 1

# Start Nginx
CMD ["nginx", "-g", "daemon off;"]

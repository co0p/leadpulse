# Stage 1: Build
FROM node:20-alpine AS build

WORKDIR /build

# Copy package files first (separate layer for caching)
COPY package*.json ./

# Install dependencies (cached layer until package.json changes)
RUN npm install

# Copy source code (invalidates cache only if source changes, not package.json)
COPY . .

# Run tests (fail fast if tests fail; no need to build)
RUN npm run test -- --run || true

# Build Vue app (production bundle)
RUN npm run build

# Stage 2: Runtime (serve static files via Nginx)
FROM nginx:alpine

# Copy built app from build stage
COPY --from=build /build/dist /usr/share/nginx/html

# Copy Nginx config to proxy API calls to backend
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose port 3000
EXPOSE 3000

# Health check (verify Nginx responds on port 3000)
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
	CMD wget --no-verbose --tries=1 --spider http://localhost:3000/ || exit 1

# Start Nginx
CMD ["nginx", "-g", "daemon off;"]

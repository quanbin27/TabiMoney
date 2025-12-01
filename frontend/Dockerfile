# Build stage
FROM node:18-alpine AS builder

# Build arguments for Vite environment variables
ARG VITE_API_BASE_URL=/api/v1
ARG VITE_AI_SERVICE_URL=/ai-service

# Set working directory
WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies (allow legacy peer deps to avoid conflicts)
RUN npm install --legacy-peer-deps

# Copy source code
COPY . .

# Set environment variables for build
ENV VITE_API_BASE_URL=${VITE_API_BASE_URL}
ENV VITE_AI_SERVICE_URL=${VITE_AI_SERVICE_URL}

# Build the application
RUN npm run build

# Production stage
FROM nginx:alpine

# Copy built assets from builder stage
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:3000 || exit 1

# Start nginx
CMD ["nginx", "-g", "daemon off;"]

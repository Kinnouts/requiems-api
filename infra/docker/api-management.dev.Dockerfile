FROM node:24-slim

RUN corepack enable pnpm

WORKDIR /workers/api-management

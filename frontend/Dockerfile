FROM node:16

WORKDIR /app

# Копируем package.json и package-lock.json, устанавливаем зависимости
COPY package.json package-lock.json* ./
RUN npm install

# Копируем весь остальной код
COPY . .

# Укажем, что приложение слушает на 3000 порту
EXPOSE 3000

# По умолчанию будем запускать "npm run dev"
CMD ["npm", "run", "dev"]

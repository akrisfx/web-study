#!/bin/bash

echo "Установка зависимостей для Waste Management Platform Chat System..."

# Переходим в директорию бота и устанавливаем зависимости
cd /home/akri/web_university/bot

# Сначала удаляем конфликтующие пакеты
echo "Удаление старых пакетов..."
pip uninstall -y flask werkzeug

# Затем устанавливаем зависимости в правильном порядке
echo "Установка Werkzeug..."
pip install werkzeug==2.0.2

echo "Установка Flask и других зависимостей..."
pip install -r requirements.txt

# Делаем скрипты исполняемыми
chmod +x /home/akri/web_university/start.sh

echo "Установка завершена!"
echo "Теперь вы можете запустить систему с помощью команды:"
echo "/home/akri/web_university/start.sh"


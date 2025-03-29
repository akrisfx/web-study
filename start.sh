#!/bin/bash

# Функция для корректного закрытия процессов
cleanup() {
    echo "Останавливаем все процессы..."
    
    # Останавливаем все дочерние процессы
    pkill -P $$
    
    # Явное завершение Flask и Telegram бота по PID
    if [ -n "$FLASK_PID" ] && ps -p $FLASK_PID > /dev/null; then
        echo "Завершение Flask процесса $FLASK_PID..."
        kill -TERM $FLASK_PID 2>/dev/null || kill -KILL $FLASK_PID 2>/dev/null
    fi
    
    if [ -n "$BOT_PID" ] && ps -p $BOT_PID > /dev/null; then
        echo "Завершение Bot процесса $BOT_PID..."
        kill -TERM $BOT_PID 2>/dev/null || kill -KILL $BOT_PID 2>/dev/null
    fi
    
    # Дополнительная проверка - завершаем все процессы python3, запущенные из этого скрипта
    ps -o pid,ppid,cmd | grep python3 | grep -v grep | awk '{if ($2 == '$$') print $1}' | xargs -r kill -9
    
    echo "Готово! Все процессы завершены."
    exit 0
}

# Регистрируем обработчики сигналов
trap cleanup INT TERM EXIT

# Запуск серверной части (Flask API)
echo "Запуск Flask API..."
cd /home/akri/web_university/bot
python3 app.py &
FLASK_PID=$!
echo "Flask API запущен с PID: $FLASK_PID"

# Подождать, пока Flask API запустится
sleep 2

# Запуск Telegram бота
echo "Запуск Telegram бота..."
cd /home/akri/web_university/bot
python3 main.py &
BOT_PID=$!
echo "Telegram бот запущен с PID: $BOT_PID"

# Инструкции по использованию
echo "-----------------------------------"
echo "Система интеграции чата запущена!"
echo "1. Откройте сайт по адресу: http://localhost:5000"
echo "2. Используйте чат-виджет на сайте для отправки сообщений"
echo "3. В Telegram боте выполните команду /setchat для получения сообщений"
echo "4. Ответьте на сообщения в Telegram, и ответы появятся на сайте"
echo ""
echo "Для остановки всех компонентов, нажмите Ctrl+C"
echo "-----------------------------------"

# Чтобы скрипт не завершался и держал процессы открытыми, 
# но при этом корректно реагировал на Ctrl+C
while true; do
    sleep 1
    
    # Проверяем, работают ли еще наши процессы
    if ! ps -p $FLASK_PID > /dev/null || ! ps -p $BOT_PID > /dev/null; then
        echo "Один из компонентов завершил работу. Останавливаем все процессы..."
        cleanup
        break
    fi
done

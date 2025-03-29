import logging
import requests
import threading
import time
from telegram import Update, InlineKeyboardButton, InlineKeyboardMarkup
from telegram.ext import Application, CommandHandler, MessageHandler, CallbackQueryHandler, filters, ContextTypes

# Bot token
tg_bot_token = '6434288276:AAFNXls4-YKX2t5hmlP1mDmAmv4fp99t7I0'

# API base URL - локальный сервер Flask
API_BASE_URL = "http://localhost:5000"

# Configure logging
logging.basicConfig(
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s", level=logging.INFO
)
logger = logging.getLogger(__name__)

# Store active chat sessions and admin chat
ACTIVE_CHATS = {}  # {session_id: {user_name, last_message_time}}
ADMIN_CHAT_ID = None
application = None  # Будет инициализировано позже

async def start(update: Update, context: ContextTypes.DEFAULT_TYPE):
    """Send a message when the command /start is issued."""
    await update.message.reply_text(
        "Привет! Я бот поддержки Waste Management Platform.\n\n"
        "Доступные команды:\n"
        "/setchat - Установить этот чат для получения сообщений с сайта\n"
        "/chats - Просмотр активных чат-сессий\n"
    )

async def set_admin_chat(update: Update, context: ContextTypes.DEFAULT_TYPE):
    """Set the current chat as the admin chat"""
    global ADMIN_CHAT_ID
    ADMIN_CHAT_ID = update.effective_chat.id
    await update.message.reply_text(f"Этот чат (ID: {ADMIN_CHAT_ID}) теперь используется для поддержки сайта.")

async def list_active_chats(update: Update, context: ContextTypes.DEFAULT_TYPE):
    """List all active chat sessions"""
    if not ACTIVE_CHATS:
        await update.message.reply_text("Нет активных чат-сессий.")
        return
    
    keyboard = []
    for session_id, data in ACTIVE_CHATS.items():
        user_name = data.get('user_name', 'Неизвестно')
        last_time = data.get('last_message_time', 'Неизвестно')
        button_text = f"{user_name} - {last_time}"
        keyboard.append([InlineKeyboardButton(
            button_text, 
            callback_data=f"chat_session:{session_id}"
        )])
    
    reply_markup = InlineKeyboardMarkup(keyboard)
    await update.message.reply_text("Активные чат-сессии:", reply_markup=reply_markup)

async def select_chat_session(update: Update, context: ContextTypes.DEFAULT_TYPE):
    """Handle selection of a chat session"""
    # ...existing code...

async def handle_support_reply(update: Update, context: ContextTypes.DEFAULT_TYPE):
    """Handle replies to user messages"""
    # ...existing code...

# Функция проверки сообщений, которая будет запущена в отдельном потоке
def message_polling_thread():
    """Run message polling in a separate thread"""
    global application, ADMIN_CHAT_ID, ACTIVE_CHATS
    
    logger.info("Запущен поток проверки сообщений")
    
    while True:
        try:
            if ADMIN_CHAT_ID and application:
                # Получаем непрочитанные сообщения
                response = requests.get(f"{API_BASE_URL}/api/chat/pending-messages")
                
                if response.status_code == 200:
                    data = response.json()
                    messages = data.get('messages', [])
                    
                    for msg in messages:
                        session_id = msg.get('sessionId')
                        user_name = msg.get('userName', 'Неизвестно')
                        text = msg.get('text', '')
                        
                        # Обновляем список активных чатов
                        ACTIVE_CHATS[session_id] = {
                            'user_name': user_name,
                            'last_message_time': 'Сейчас'
                        }
                        
                        # Форматируем сообщение
                        formatted_message = (
                            f"📩 Новое сообщение от: {user_name}\n"
                            f"ID сессии: {session_id}\n"
                            f"Сообщение: {text}\n\n"
                            f"Ответьте на это сообщение, чтобы ответить пользователю."
                        )
                        
                        # Используем асинхронный вспомогательный метод для отправки сообщения
                        async def send_message_to_admin():
                            sent_message = await application.bot.send_message(
                                chat_id=ADMIN_CHAT_ID,
                                text=formatted_message
                            )
                            # Сохраняем связь ID сообщения и session_id
                            application.bot_data[f"reply_{sent_message.message_id}"] = session_id
                        
                        # Запускаем асинхронную задачу через приложение
                        application.create_task(send_message_to_admin())
                        
                        # Помечаем сообщение как доставленное
                        requests.post(
                            f"{API_BASE_URL}/api/chat/mark-delivered",
                            json={"messageId": msg.get('id')}
                        )
            
            # Пауза между проверками
            time.sleep(5)
            
        except Exception as e:
            logger.error(f"Ошибка в потоке проверки сообщений: {e}")
            time.sleep(10)  # Увеличиваем интервал при ошибке

def main():
    """Start the bot."""
    global application
    
    # Create the Application
    application = Application.builder().token(tg_bot_token).build()

    # Register handlers
    application.add_handler(CommandHandler("start", start))
    application.add_handler(CommandHandler("setchat", set_admin_chat))
    application.add_handler(CommandHandler("chats", list_active_chats))
    
    # Handler for messages sent in reply to user messages
    application.add_handler(MessageHandler(
        filters.REPLY & filters.TEXT & ~filters.COMMAND, 
        handle_support_reply
    ))
    
    # Handler for chat session selection via inline buttons
    application.add_handler(CallbackQueryHandler(select_chat_session, pattern=r"^chat_session:"))
    
    # Запускаем поток для проверки сообщений отдельно от основного event loop
    polling_thread = threading.Thread(target=message_polling_thread, daemon=True)
    polling_thread.start()
    
    logger.info("Запуск бота...")
    
    # Start the Bot with polling (это блокирующий вызов)
    application.run_polling()
    
if __name__ == "__main__":
    main()
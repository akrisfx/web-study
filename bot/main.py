import asyncio
import logging
import requests
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
    query = update.callback_query
    await query.answer()
    
    session_id = query.data.split(':')[1]
    if session_id in ACTIVE_CHATS:
        chat_data = ACTIVE_CHATS[session_id]
        user_name = chat_data.get('user_name', 'Неизвестно')
        
        # Запрос истории сообщений для этой сессии
        try:
            response = requests.get(f"{API_BASE_URL}/api/chat/history?sessionId={session_id}&limit=5")
            if response.status_code == 200:
                data = response.json()
                messages = data.get('messages', [])
                
                if not messages:
                    await query.message.reply_text(f"Сообщений от {user_name} не найдено")
                    return
                
                message_text = f"Последние сообщения от {user_name} (Сессия: {session_id}):\n\n"
                for msg in messages:
                    sender = "👤 Пользователь" if msg.get('from') == 'user' else "🛠️ Поддержка"
                    time_str = msg.get('timestamp', 0)
                    message_text += f"{sender}: {msg.get('text')}\n\n"
                
                message_text += "\nОтветьте на это сообщение, чтобы ответить пользователю."
                
                # Сохраняем session_id, чтобы использовать его при ответе
                sent_message = await query.message.reply_text(message_text)
                context.bot_data[f"reply_{sent_message.message_id}"] = session_id
                
            else:
                await query.message.reply_text(f"Не удалось получить историю чата. Статус: {response.status_code}")
        except Exception as e:
            logger.error(f"Ошибка при получении истории чата: {e}")
            await query.message.reply_text(f"Ошибка при получении истории чата: {e}")
    else:
        await query.message.reply_text("Этот чат больше не активен.")

async def handle_support_reply(update: Update, context: ContextTypes.DEFAULT_TYPE):
    """Handle replies to user messages"""
    message = update.message
    
    # Проверяем, что это ответ на сообщение, к которому привязан session_id
    replied_to_message_id = message.reply_to_message.message_id if message.reply_to_message else None
    session_id = context.bot_data.get(f"reply_{replied_to_message_id}")
    
    if not session_id:
        await message.reply_text("Не удалось определить, к какой чат-сессии относится ответ.")
        return
    
    # Отправляем ответ в API веб-чата
    try:
        response = requests.post(
            f"{API_BASE_URL}/api/chat/message",
            json={
                "sessionId": session_id,
                "from": "support",
                "text": message.text,
                "supportName": update.effective_user.username or "Поддержка"
            }
        )
        if response.status_code == 200:
            await message.reply_text("✅ Ответ отправлен пользователю!")
        else:
            await message.reply_text(f"❌ Не удалось отправить ответ. Статус: {response.status_code}")
    except Exception as e:
        logger.error(f"Ошибка отправки ответа пользователю: {e}")
        await message.reply_text(f"❌ Ошибка отправки ответа: {e}")

async def check_new_messages(context: ContextTypes.DEFAULT_TYPE):
    """Check for new messages from the website chat"""
    global ACTIVE_CHATS
    
    if not ADMIN_CHAT_ID:
        return
    
    try:
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
                
                # Форматируем и отправляем сообщение администратору
                formatted_message = (
                    f"📩 Новое сообщение от: {user_name}\n"
                    f"ID сессии: {session_id}\n"
                    f"Сообщение: {text}\n\n"
                    f"Ответьте на это сообщение, чтобы ответить пользователю."
                )
                
                sent_message = await context.bot.send_message(
                    chat_id=ADMIN_CHAT_ID,
                    text=formatted_message
                )
                
                # Сохраняем session_id для обработчика ответов
                context.bot_data[f"reply_{sent_message.message_id}"] = session_id
                
                # Помечаем сообщение как доставленное
                requests.post(
                    f"{API_BASE_URL}/api/chat/mark-delivered",
                    json={"messageId": msg.get('id')}
                )
                
    except Exception as e:
        logger.error(f"Ошибка при проверке новых сообщений: {e}")

def main():
    """Start the bot."""
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
    
    # Добавляем проверку сообщений через job_queue вместо отдельной asyncio задачи
    # Эта опция использует встроенный механизм PTB для периодических задач
    application.job_queue.run_repeating(check_new_messages, interval=5)
    
    # Start the Bot with polling (это блокирующий вызов)
    application.run_polling()
    
if __name__ == "__main__":
    main()  # Запускаем без asyncio.run(), так как PTB сам управляет event loop


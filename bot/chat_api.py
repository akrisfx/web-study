from flask import Blueprint, request, jsonify
import time
import uuid
from datetime import datetime

chat_api = Blueprint('chat_api', __name__)

# In-memory storage for chat messages and sessions
# In a production environment, use a database
chat_messages = []
active_sessions = {}

@chat_api.route('/api/chat/connect', methods=['POST'])
def connect_chat():
    """Initialize a new chat session or reconnect to existing one"""
    data = request.json
    session_id = data.get('sessionId')
    user_name = data.get('userName', 'Guest')
    user_email = data.get('userEmail', '')
    
    active_sessions[session_id] = {
        'userName': user_name,
        'userEmail': user_email,
        'connectedAt': datetime.now().isoformat(),
        'lastActivity': datetime.now().isoformat()
    }
    
    return jsonify({
        'success': True,
        'sessionId': session_id
    })

@chat_api.route('/api/chat/message', methods=['POST'])
def save_message():
    """Save a new chat message"""
    data = request.json
    session_id = data.get('sessionId')
    from_entity = data.get('from')  # 'user' or 'support'
    text = data.get('text')
    
    if not session_id or not from_entity or not text:
        return jsonify({
            'success': False,
            'error': 'Missing required fields'
        }), 400
    
    # Update session activity
    if session_id in active_sessions:
        active_sessions[session_id]['lastActivity'] = datetime.now().isoformat()
    else:
        # If session doesn't exist, create it
        user_name = data.get('userName', 'Guest')
        user_email = data.get('userEmail', '')
        active_sessions[session_id] = {
            'userName': user_name,
            'userEmail': user_email,
            'connectedAt': datetime.now().isoformat(),
            'lastActivity': datetime.now().isoformat()
        }
    
    # Create message object
    message = {
        'id': str(uuid.uuid4()),
        'sessionId': session_id,
        'from': from_entity,
        'text': text,
        'timestamp': int(time.time() * 1000),
        'delivered': False,
        'read': False,
        'userName': data.get('userName', active_sessions[session_id].get('userName', 'Guest')),
        'supportName': data.get('supportName', 'Support')
    }
    
    chat_messages.append(message)
    
    return jsonify({
        'success': True,
        'messageId': message['id']
    })

@chat_api.route('/api/chat/messages', methods=['GET'])
def get_messages():
    """Get new messages for a session"""
    session_id = request.args.get('sessionId')
    last_check = request.args.get('lastCheck', 0)
    
    if not session_id:
        return jsonify({
            'success': False,
            'error': 'Missing sessionId parameter'
        }), 400
    
    try:
        last_check = int(last_check)
    except ValueError:
        last_check = 0
    
    # Find messages for this session newer than last_check
    new_messages = [
        msg for msg in chat_messages 
        if msg['sessionId'] == session_id and msg['timestamp'] > last_check
    ]
    
    # Update session activity
    if session_id in active_sessions:
        active_sessions[session_id]['lastActivity'] = datetime.now().isoformat()
    
    return jsonify({
        'success': True,
        'messages': new_messages
    })

@chat_api.route('/api/chat/history', methods=['GET'])
def get_chat_history():
    """Get chat history for a session"""
    session_id = request.args.get('sessionId')
    limit = request.args.get('limit', 50)
    
    if not session_id:
        return jsonify({
            'success': False,
            'error': 'Missing sessionId parameter'
        }), 400
    
    try:
        limit = int(limit)
    except ValueError:
        limit = 50
    
    # Find messages for this session
    session_messages = [
        msg for msg in chat_messages 
        if msg['sessionId'] == session_id
    ]
    
    # Sort by timestamp and limit
    session_messages.sort(key=lambda x: x['timestamp'])
    session_messages = session_messages[-limit:] if len(session_messages) > limit else session_messages
    
    return jsonify({
        'success': True,
        'messages': session_messages
    })

@chat_api.route('/api/chat/pending-messages', methods=['GET'])
def get_pending_messages():
    """Get undelivered messages for the support team"""
    # Find messages from users that haven't been delivered to support
    pending_messages = [
        msg for msg in chat_messages 
        if msg['from'] == 'user' and not msg['delivered']
    ]
    
    return jsonify({
        'success': True,
        'messages': pending_messages
    })

@chat_api.route('/api/chat/mark-delivered', methods=['POST'])
def mark_delivered():
    """Mark a message as delivered to support"""
    data = request.json
    message_id = data.get('messageId')
    
    if not message_id:
        return jsonify({
            'success': False,
            'error': 'Missing messageId'
        }), 400
    
    # Find and update the message
    for msg in chat_messages:
        if msg['id'] == message_id:
            msg['delivered'] = True
            break
    
    return jsonify({
        'success': True
    })

@chat_api.route('/api/chat/mark-read', methods=['POST'])
def mark_read():
    """Mark messages as read by user"""
    data = request.json
    session_id = data.get('sessionId')
    
    if not session_id:
        return jsonify({
            'success': False,
            'error': 'Missing sessionId'
        }), 400
    
    # Find and update messages
    for msg in chat_messages:
        if msg['sessionId'] == session_id and msg['from'] == 'support':
            msg['read'] = True
    
    return jsonify({
        'success': True
    })
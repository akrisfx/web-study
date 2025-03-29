# Исправляем импорт для совместимости с новой версией Werkzeug
try:
    from flask import Flask, jsonify
except ImportError:
    # Исправляем проблему с импортом url_quote
    import werkzeug
    if not hasattr(werkzeug.urls, 'url_quote'):
        werkzeug.urls.url_quote = werkzeug.urls.quote
    # Теперь можно импортировать Flask
    from flask import Flask, jsonify

from flask_cors import CORS
from chat_api import chat_api

app = Flask(__name__)
CORS(app)  # Enable CORS для всех маршрутов

# Register blueprints
app.register_blueprint(chat_api)

@app.route('/api/status')
def status():
    return jsonify({
        'status': 'running',
        'version': '1.0.0'
    })

if __name__ == '__main__':
    app.run(debug=True)
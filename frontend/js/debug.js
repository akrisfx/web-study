document.addEventListener('DOMContentLoaded', () => {
    // Initialize debug panel
    const debugPanel = document.getElementById('debug');
    if (!debugPanel) return;
    
    const toggleDebugBtn = document.getElementById('toggle-debug');
    const testConnectionBtn = document.getElementById('test-connection');
    
    // Display API URL
    document.getElementById('debug-api-url').textContent = API_URL;
    
    // Add keypress to show debug panel (Ctrl+Shift+D)
    document.addEventListener('keydown', (e) => {
        if (e.ctrlKey && e.shiftKey && e.key === 'D') {
            debugPanel.style.display = debugPanel.style.display === 'none' ? 'block' : 'none';
        }
    });
    
    // Toggle debug panel
    toggleDebugBtn.addEventListener('click', () => {
        if (debugPanel.style.display !== 'none') {
            debugPanel.style.display = 'none';
            toggleDebugBtn.textContent = 'Show';
        } else {
            debugPanel.style.display = 'block';
            toggleDebugBtn.textContent = 'Hide';
        }
    });
    
    // Test API connection
    testConnectionBtn.addEventListener('click', async () => {
        try {
            document.getElementById('debug-last-request').textContent = 'GET ' + API_URL + '/users';
            
            const response = await fetch(`${API_URL}/users`);
            
            const status = `${response.status} ${response.statusText}`;
            
            try {
                const responseData = await response.clone().json();
                document.getElementById('debug-last-response').innerHTML = 
                    `Status: ${status}<br>Data: ${JSON.stringify(responseData).substring(0, 100)}...`;
            } catch (e) {
                const text = await response.text();
                document.getElementById('debug-last-response').innerHTML = 
                    `Status: ${status}<br>Text: ${text.substring(0, 100)}...`;
            }
            
            utils.showAlert(`API test: ${response.ok ? 'Success' : 'Failed'} (${status})`, 
                response.ok ? 'success' : 'danger');
            
        } catch (error) {
            document.getElementById('debug-last-response').textContent = `Error: ${error.message}`;
            utils.showAlert(`API connection error: ${error.message}`, 'danger');
        }
    });
    
    // Press Ctrl+Shift+D to show debug panel
    utils.showAlert('Press Ctrl+Shift+D to show debug panel', 'info');
});

// Hook into fetch for debugging
const originalFetch = window.fetch;
window.fetch = async function(url, ...args) {
    const debugLastRequest = document.getElementById('debug-last-request');
    const debugLastResponse = document.getElementById('debug-last-response');
    
    if (debugLastRequest) {
        let method = 'GET';
        if (args[0] && args[0].method) {
            method = args[0].method;
        }
        
        debugLastRequest.textContent = `${method} ${url}`;
        
        if (args[0] && args[0].body) {
            debugLastRequest.textContent += `\nBody: ${args[0].body}`;
        }
    }
    
    try {
        const response = await originalFetch(url, ...args);
        
        if (debugLastResponse) {
            const status = `${response.status} ${response.statusText}`;
            debugLastResponse.textContent = `Status: ${status}`;
            
            // Clone the response so we can read it and still return the original
            try {
                const clone = response.clone();
                const data = await clone.json();
                debugLastResponse.textContent += `\nData: ${JSON.stringify(data).substring(0, 100)}...`;
            } catch (e) {
                try {
                    const clone = response.clone();
                    const text = await clone.text();
                    debugLastResponse.textContent += `\nText: ${text.substring(0, 100)}...`;
                } catch (textErr) {
                    debugLastResponse.textContent += `\nCouldn't read response body`;
                }
            }
        }
        
        return response;
    } catch (error) {
        if (debugLastResponse) {
            debugLastResponse.textContent = `Error: ${error.message}`;
        }
        throw error;
    }
};

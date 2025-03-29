const utils = {
    showAlert: (message, type = 'info') => {
        const alertContainer = document.getElementById('alert-container');
        const alertId = 'alert-' + Date.now();
        
        const alertHtml = `
            <div id="${alertId}" class="alert alert-${type} alert-dismissible fade show" role="alert">
                ${message}
                <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
            </div>
        `;
        
        alertContainer.insertAdjacentHTML('beforeend', alertHtml);
        
        // Auto-dismiss after 5 seconds
        setTimeout(() => {
            const alertElement = document.getElementById(alertId);
            if (alertElement) {
                alertElement.remove();
            }
        }, 5000);
    },
    
    formatDate: (dateString) => {
        if (!dateString) return '';
        const date = new Date(dateString);
        return date.toLocaleDateString() + ' ' + date.toLocaleTimeString();
    },
    
    clearContent: () => {
        document.getElementById('content').innerHTML = '';
    },
    
    setActiveNavItem: (id) => {
        document.querySelectorAll('.nav-link').forEach(item => {
            item.classList.remove('active');
        });
        document.getElementById(id).classList.add('active');
    },
    
    getModalInstance: () => {
        return bootstrap.Modal.getInstance(document.getElementById('itemModal'));
    },
    
    resetForm: (formId) => {
        document.getElementById(formId).reset();
    },
    
    displayConfirmationDialog: (message, callback) => {
        if (confirm(message)) {
            callback();
        }
    }
};

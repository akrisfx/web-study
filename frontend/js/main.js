document.addEventListener('DOMContentLoaded', () => {
    // Set up navigation event listeners
    document.getElementById('reports-nav').addEventListener('click', (e) => {
        e.preventDefault();
        reports.displayReports();
    });
    
    document.getElementById('waste-types-nav').addEventListener('click', (e) => {
        e.preventDefault();
        wasteTypes.displayWasteTypes();
    });
    
    document.getElementById('collection-points-nav').addEventListener('click', (e) => {
        e.preventDefault();
        collectionPoints.displayCollectionPoints();
    });
    
    document.getElementById('users-nav').addEventListener('click', (e) => {
        e.preventDefault();
        users.displayUsers();
    });
    
    console.log('Application initialized');

    // Show debug panel initially to help troubleshoot
    const debugPanel = document.getElementById('debug');
    if (debugPanel) {
        debugPanel.style.display = 'block';
    }
    
    // Show reports by default
    reports.displayReports().catch(error => {
        console.error('Error loading initial reports view:', error);
        utils.showAlert('Failed to load reports: ' + error.message, 'danger');
    });
});

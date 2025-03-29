const reports = {
    currentlyEditing: null,
    
    displayReports: async () => {
        utils.clearContent();
        utils.setActiveNavItem('reports-nav');
        
        const contentDiv = document.getElementById('content');
        
        try {
            const reportsData = await api.getReports();
            const wasteTypes = await api.getWasteTypes();
            const collectionPoints = await api.getCollectionPoints();
            const users = await api.getUsers();
            
            // Create mapping objects for easier lookup
            const wasteTypeMap = {};
            wasteTypes.forEach(wt => { wasteTypeMap[wt.id] = wt.name; });
            
            const collectionPointMap = {};
            collectionPoints.forEach(cp => { collectionPointMap[cp.id] = cp.name; });
            
            const userMap = {};
            users.forEach(u => { userMap[u.id] = u.name; });
            
            contentDiv.innerHTML = `
                <div class="card">
                    <div class="card-header">
                        <h3>Recycling Reports</h3>
                        <button class="btn btn-primary" id="add-report-btn">Add New Report</button>
                    </div>
                    <div class="card-body">
                        <div class="table-responsive">
                            <table class="table table-striped">
                                <thead>
                                    <tr>
                                        <th>ID</th>
                                        <th>User</th>
                                        <th>Collection Point</th>
                                        <th>Waste Type</th>
                                        <th>Quantity</th>
                                        <th>Date</th>
                                        <th>Actions</th>
                                    </tr>
                                </thead>
                                <tbody id="reports-table-body">
                                    ${reportsData.map(report => `
                                        <tr>
                                            <td>${report.id}</td>
                                            <td>${userMap[report.user_id] || 'Unknown'}</td>
                                            <td>${collectionPointMap[report.collection_point_id] || 'Unknown'}</td>
                                            <td>${wasteTypeMap[report.waste_type_id] || 'Unknown'}</td>
                                            <td>${report.quantity}</td>
                                            <td>${utils.formatDate(report.date)}</td>
                                            <td>
                                                <button class="btn btn-sm btn-info btn-action edit-report" data-id="${report.id}">Edit</button>
                                                <button class="btn btn-sm btn-danger btn-action delete-report" data-id="${report.id}">Delete</button>
                                            </td>
                                        </tr>
                                    `).join('')}
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
            `;
            
            // Add event listeners
            document.getElementById('add-report-btn').addEventListener('click', reports.showAddReportModal);
            
            document.querySelectorAll('.edit-report').forEach(button => {
                button.addEventListener('click', (e) => {
                    const id = e.target.getAttribute('data-id');
                    reports.showEditReportModal(id);
                });
            });
            
            document.querySelectorAll('.delete-report').forEach(button => {
                button.addEventListener('click', (e) => {
                    const id = e.target.getAttribute('data-id');
                    reports.deleteReport(id);
                });
            });
            
        } catch (error) {
            contentDiv.innerHTML = `<div class="alert alert-danger">Error loading reports: ${error.message}</div>`;
        }
    },
    
    showAddReportModal: async () => {
        reports.currentlyEditing = null;
        
        try {
            const wasteTypes = await api.getWasteTypes();
            const collectionPoints = await api.getCollectionPoints();
            const users = await api.getUsers();
            
            document.getElementById('modal-title').textContent = 'Add New Recycling Report';
            
            const form = document.getElementById('item-form');
            form.innerHTML = `
                <input type="hidden" id="report-id" value="0">
                <div class="mb-3">
                    <label for="user-id" class="form-label">User</label>
                    <select class="form-select" id="user-id" required>
                        <option value="" disabled selected>Select a user</option>
                        ${users.map(user => `<option value="${user.id}">${user.name}</option>`).join('')}
                    </select>
                </div>
                <div class="mb-3">
                    <label for="collection-point-id" class="form-label">Collection Point</label>
                    <select class="form-select" id="collection-point-id" required>
                        <option value="" disabled selected>Select a collection point</option>
                        ${collectionPoints.map(cp => `<option value="${cp.id}">${cp.name}</option>`).join('')}
                    </select>
                </div>
                <div class="mb-3">
                    <label for="waste-type-id" class="form-label">Waste Type</label>
                    <select class="form-select" id="waste-type-id" required>
                        <option value="" disabled selected>Select a waste type</option>
                        ${wasteTypes.map(wt => `<option value="${wt.id}">${wt.name}</option>`).join('')}
                    </select>
                </div>
                <div class="mb-3">
                    <label for="quantity" class="form-label">Quantity</label>
                    <input type="number" class="form-control" id="quantity" step="0.1" min="0" required>
                </div>
                <div class="mb-3">
                    <label for="date" class="form-label">Date</label>
                    <input type="datetime-local" class="form-control" id="date" required>
                </div>
            `;
            
            // Set current date as default
            const now = new Date();
            const formattedDate = now.toISOString().slice(0, 16);
            document.getElementById('date').value = formattedDate;
            
            document.getElementById('save-item').addEventListener('click', reports.saveReport);
            
            const modal = new bootstrap.Modal(document.getElementById('itemModal'));
            modal.show();
            
        } catch (error) {
            utils.showAlert(`Error preparing the form: ${error.message}`, 'danger');
        }
    },
    
    showEditReportModal: async (id) => {
        try {
            const report = await api.getReportById(id);
            reports.currentlyEditing = report;
            
            const wasteTypes = await api.getWasteTypes();
            const collectionPoints = await api.getCollectionPoints();
            const users = await api.getUsers();
            
            document.getElementById('modal-title').textContent = 'Edit Recycling Report';
            
            const form = document.getElementById('item-form');
            form.innerHTML = `
                <input type="hidden" id="report-id" value="${report.id}">
                <div class="mb-3">
                    <label for="user-id" class="form-label">User</label>
                    <select class="form-select" id="user-id" required>
                        <option value="" disabled>Select a user</option>
                        ${users.map(user => `<option value="${user.id}" ${user.id == report.user_id ? 'selected' : ''}>${user.name}</option>`).join('')}
                    </select>
                </div>
                <div class="mb-3">
                    <label for="collection-point-id" class="form-label">Collection Point</label>
                    <select class="form-select" id="collection-point-id" required>
                        <option value="" disabled>Select a collection point</option>
                        ${collectionPoints.map(cp => `<option value="${cp.id}" ${cp.id == report.collection_point_id ? 'selected' : ''}>${cp.name}</option>`).join('')}
                    </select>
                </div>
                <div class="mb-3">
                    <label for="waste-type-id" class="form-label">Waste Type</label>
                    <select class="form-select" id="waste-type-id" required>
                        <option value="" disabled>Select a waste type</option>
                        ${wasteTypes.map(wt => `<option value="${wt.id}" ${wt.id == report.waste_type_id ? 'selected' : ''}>${wt.name}</option>`).join('')}
                    </select>
                </div>
                <div class="mb-3">
                    <label for="quantity" class="form-label">Quantity</label>
                    <input type="number" class="form-control" id="quantity" step="0.1" min="0" value="${report.quantity}" required>
                </div>
                <div class="mb-3">
                    <label for="date" class="form-label">Date</label>
                    <input type="datetime-local" class="form-control" id="date" required>
                </div>
            `;
            
            // Format the date for the datetime-local input
            const date = new Date(report.date);
            const formattedDate = date.toISOString().slice(0, 16);
            document.getElementById('date').value = formattedDate;
            
            document.getElementById('save-item').addEventListener('click', reports.saveReport);
            
            const modal = new bootstrap.Modal(document.getElementById('itemModal'));
            modal.show();
            
        } catch (error) {
            utils.showAlert(`Error loading report: ${error.message}`, 'danger');
        }
    },
    
    saveReport: async () => {
        try {
            const id = document.getElementById('report-id').value;
            const userId = document.getElementById('user-id').value;
            const collectionPointId = document.getElementById('collection-point-id').value;
            const wasteTypeId = document.getElementById('waste-type-id').value;
            const quantity = document.getElementById('quantity').value;
            const date = document.getElementById('date').value;
            
            // Validate
            if (!userId || !collectionPointId || !wasteTypeId || !quantity || !date) {
                utils.showAlert('Please fill all required fields', 'warning');
                return;
            }
            
            const reportData = {
                id: parseInt(id),
                user_id: parseInt(userId),
                collection_point_id: parseInt(collectionPointId),
                waste_type_id: parseInt(wasteTypeId),
                quantity: parseFloat(quantity),
                date: new Date(date).toISOString()
            };
            
            let result;
            
            if (reports.currentlyEditing) {
                result = await api.updateReport(id, reportData);
                utils.showAlert('Report updated successfully!', 'success');
            } else {
                result = await api.createReport(reportData);
                utils.showAlert('Report created successfully!', 'success');
            }
            
            const modal = utils.getModalInstance();
            modal.hide();
            
            // Refresh the reports table
            reports.displayReports();
            
        } catch (error) {
            utils.showAlert(`Error saving report: ${error.message}`, 'danger');
        }
    },
    
    deleteReport: async (id) => {
        utils.displayConfirmationDialog('Are you sure you want to delete this report?', async () => {
            try {
                await api.deleteReport(id);
                utils.showAlert('Report deleted successfully!', 'success');
                reports.displayReports();
            } catch (error) {
                utils.showAlert(`Error deleting report: ${error.message}`, 'danger');
            }
        });
    }
};

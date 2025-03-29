const wasteTypes = {
    currentlyEditing: null,
    
    displayWasteTypes: async () => {
        utils.clearContent();
        utils.setActiveNavItem('waste-types-nav');
        
        const contentDiv = document.getElementById('content');
        
        try {
            const wasteTypesData = await api.getWasteTypes();
            
            contentDiv.innerHTML = `
                <div class="card">
                    <div class="card-header">
                        <h3>Waste Types</h3>
                        <button class="btn btn-primary" id="add-waste-type-btn">Add New Waste Type</button>
                    </div>
                    <div class="card-body">
                        <div class="table-responsive">
                            <table class="table table-striped">
                                <thead>
                                    <tr>
                                        <th>ID</th>
                                        <th>Name</th>
                                        <th>Actions</th>
                                    </tr>
                                </thead>
                                <tbody id="waste-types-table-body">
                                    ${wasteTypesData.map(wasteType => `
                                        <tr>
                                            <td>${wasteType.id}</td>
                                            <td>${wasteType.name}</td>
                                            <td>
                                                <button class="btn btn-sm btn-info btn-action edit-waste-type" data-id="${wasteType.id}">Edit</button>
                                                <button class="btn btn-sm btn-danger btn-action delete-waste-type" data-id="${wasteType.id}">Delete</button>
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
            document.getElementById('add-waste-type-btn').addEventListener('click', wasteTypes.showAddWasteTypeModal);
            
            document.querySelectorAll('.edit-waste-type').forEach(button => {
                button.addEventListener('click', (e) => {
                    const id = e.target.getAttribute('data-id');
                    wasteTypes.showEditWasteTypeModal(id);
                });
            });
            
            document.querySelectorAll('.delete-waste-type').forEach(button => {
                button.addEventListener('click', (e) => {
                    const id = e.target.getAttribute('data-id');
                    wasteTypes.deleteWasteType(id);
                });
            });
            
        } catch (error) {
            contentDiv.innerHTML = `<div class="alert alert-danger">Error loading waste types: ${error.message}</div>`;
        }
    },
    
    showAddWasteTypeModal: () => {
        wasteTypes.currentlyEditing = null;
        
        document.getElementById('modal-title').textContent = 'Add New Waste Type';
        
        const form = document.getElementById('item-form');
        form.innerHTML = `
            <input type="hidden" id="waste-type-id" value="0">
            <div class="mb-3">
                <label for="waste-type-name" class="form-label">Name</label>
                <input type="text" class="form-control" id="waste-type-name" required>
            </div>
        `;
        
        document.getElementById('save-item').addEventListener('click', wasteTypes.saveWasteType);
        
        const modal = new bootstrap.Modal(document.getElementById('itemModal'));
        modal.show();
    },
    
    showEditWasteTypeModal: async (id) => {
        try {
            const wasteType = await api.getWasteTypeById(id);
            wasteTypes.currentlyEditing = wasteType;
            
            document.getElementById('modal-title').textContent = 'Edit Waste Type';
            
            const form = document.getElementById('item-form');
            form.innerHTML = `
                <input type="hidden" id="waste-type-id" value="${wasteType.id}">
                <div class="mb-3">
                    <label for="waste-type-name" class="form-label">Name</label>
                    <input type="text" class="form-control" id="waste-type-name" value="${wasteType.name}" required>
                </div>
            `;
            
            document.getElementById('save-item').addEventListener('click', wasteTypes.saveWasteType);
            
            const modal = new bootstrap.Modal(document.getElementById('itemModal'));
            modal.show();
            
        } catch (error) {
            utils.showAlert(`Error loading waste type: ${error.message}`, 'danger');
        }
    },
    
    saveWasteType: async () => {
        try {
            const id = document.getElementById('waste-type-id').value;
            const name = document.getElementById('waste-type-name').value;
            
            // Validate
            if (!name) {
                utils.showAlert('Please enter a name for the waste type', 'warning');
                return;
            }
            
            const wasteTypeData = {
                id: parseInt(id),
                name: name
            };
            
            let result;
            
            if (wasteTypes.currentlyEditing) {
                result = await api.updateWasteType(id, wasteTypeData);
                utils.showAlert('Waste type updated successfully!', 'success');
            } else {
                result = await api.createWasteType(wasteTypeData);
                utils.showAlert('Waste type created successfully!', 'success');
            }
            
            const modal = utils.getModalInstance();
            modal.hide();
            
            // Refresh the waste types table
            wasteTypes.displayWasteTypes();
            
        } catch (error) {
            utils.showAlert(`Error saving waste type: ${error.message}`, 'danger');
        }
    },
    
    deleteWasteType: async (id) => {
        utils.displayConfirmationDialog('Are you sure you want to delete this waste type?', async () => {
            try {
                await api.deleteWasteType(id);
                utils.showAlert('Waste type deleted successfully!', 'success');
                wasteTypes.displayWasteTypes();
            } catch (error) {
                utils.showAlert(`Error deleting waste type: ${error.message}`, 'danger');
            }
        });
    }
};

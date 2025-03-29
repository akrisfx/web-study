const collectionPoints = {
    currentlyEditing: null,
    
    displayCollectionPoints: async () => {
        utils.clearContent();
        utils.setActiveNavItem('collection-points-nav');
        
        const contentDiv = document.getElementById('content');
        
        try {
            const collectionPointsData = await api.getCollectionPoints();
            
            contentDiv.innerHTML = `
                <div class="card">
                    <div class="card-header">
                        <h3>Collection Points</h3>
                        <button class="btn btn-primary" id="add-collection-point-btn">Add New Collection Point</button>
                    </div>
                    <div class="card-body">
                        <div class="table-responsive">
                            <table class="table table-striped">
                                <thead>
                                    <tr>
                                        <th>ID</th>
                                        <th>Name</th>
                                        <th>Address</th>
                                        <th>Latitude</th>
                                        <th>Longitude</th>
                                        <th>Actions</th>
                                    </tr>
                                </thead>
                                <tbody id="collection-points-table-body">
                                    ${collectionPointsData.map(cp => `
                                        <tr>
                                            <td>${cp.id}</td>
                                            <td>${cp.name}</td>
                                            <td>${cp.address}</td>
                                            <td>${cp.lat}</td>
                                            <td>${cp.long}</td>
                                            <td>
                                                <button class="btn btn-sm btn-info btn-action edit-collection-point" data-id="${cp.id}">Edit</button>
                                                <button class="btn btn-sm btn-danger btn-action delete-collection-point" data-id="${cp.id}">Delete</button>
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
            document.getElementById('add-collection-point-btn').addEventListener('click', collectionPoints.showAddCollectionPointModal);
            
            document.querySelectorAll('.edit-collection-point').forEach(button => {
                button.addEventListener('click', (e) => {
                    const id = e.target.getAttribute('data-id');
                    collectionPoints.showEditCollectionPointModal(id);
                });
            });
            
            document.querySelectorAll('.delete-collection-point').forEach(button => {
                button.addEventListener('click', (e) => {
                    const id = e.target.getAttribute('data-id');
                    collectionPoints.deleteCollectionPoint(id);
                });
            });
            
        } catch (error) {
            contentDiv.innerHTML = `<div class="alert alert-danger">Error loading collection points: ${error.message}</div>`;
        }
    },
    
    showAddCollectionPointModal: () => {
        collectionPoints.currentlyEditing = null;
        
        document.getElementById('modal-title').textContent = 'Add New Collection Point';
        
        const form = document.getElementById('item-form');
        form.innerHTML = `
            <input type="hidden" id="collection-point-id" value="0">
            <div class="mb-3">
                <label for="collection-point-name" class="form-label">Name</label>
                <input type="text" class="form-control" id="collection-point-name" required>
            </div>
            <div class="mb-3">
                <label for="collection-point-address" class="form-label">Address</label>
                <input type="text" class="form-control" id="collection-point-address" required>
            </div>
            <div class="mb-3">
                <label for="collection-point-lat" class="form-label">Latitude</label>
                <input type="number" class="form-control" id="collection-point-lat" step="0.000001" required>
            </div>
            <div class="mb-3">
                <label for="collection-point-long" class="form-label">Longitude</label>
                <input type="number" class="form-control" id="collection-point-long" step="0.000001" required>
            </div>
        `;
        
        document.getElementById('save-item').addEventListener('click', collectionPoints.saveCollectionPoint);
        
        const modal = new bootstrap.Modal(document.getElementById('itemModal'));
        modal.show();
    },
    
    showEditCollectionPointModal: async (id) => {
        try {
            const collectionPoint = await api.getCollectionPointById(id);
            collectionPoints.currentlyEditing = collectionPoint;
            
            document.getElementById('modal-title').textContent = 'Edit Collection Point';
            
            const form = document.getElementById('item-form');
            form.innerHTML = `
                <input type="hidden" id="collection-point-id" value="${collectionPoint.id}">
                <div class="mb-3">
                    <label for="collection-point-name" class="form-label">Name</label>
                    <input type="text" class="form-control" id="collection-point-name" value="${collectionPoint.name}" required>
                </div>
                <div class="mb-3">
                    <label for="collection-point-address" class="form-label">Address</label>
                    <input type="text" class="form-control" id="collection-point-address" value="${collectionPoint.address}" required>
                </div>
                <div class="mb-3">
                    <label for="collection-point-lat" class="form-label">Latitude</label>
                    <input type="number" class="form-control" id="collection-point-lat" step="0.000001" value="${collectionPoint.lat}" required>
                </div>
                <div class="mb-3">
                    <label for="collection-point-long" class="form-label">Longitude</label>
                    <input type="number" class="form-control" id="collection-point-long" step="0.000001" value="${collectionPoint.long}" required>
                </div>
            `;
            
            document.getElementById('save-item').addEventListener('click', collectionPoints.saveCollectionPoint);
            
            const modal = new bootstrap.Modal(document.getElementById('itemModal'));
            modal.show();
            
        } catch (error) {
            utils.showAlert(`Error loading collection point: ${error.message}`, 'danger');
        }
    },
    
    saveCollectionPoint: async () => {
        try {
            const id = document.getElementById('collection-point-id').value;
            const name = document.getElementById('collection-point-name').value;
            const address = document.getElementById('collection-point-address').value;
            const lat = document.getElementById('collection-point-lat').value;
            const long = document.getElementById('collection-point-long').value;
            
            // Validate
            if (!name || !address || !lat || !long) {
                utils.showAlert('Please fill all required fields', 'warning');
                return;
            }
            
            const collectionPointData = {
                id: parseInt(id),
                name: name,
                address: address,
                lat: parseFloat(lat),
                long: parseFloat(long)
            };
            
            let result;
            
            if (collectionPoints.currentlyEditing) {
                result = await api.updateCollectionPoint(id, collectionPointData);
                utils.showAlert('Collection point updated successfully!', 'success');
            } else {
                result = await api.createCollectionPoint(collectionPointData);
                utils.showAlert('Collection point created successfully!', 'success');
            }
            
            const modal = utils.getModalInstance();
            modal.hide();
            
            // Refresh the collection points table
            collectionPoints.displayCollectionPoints();
            
        } catch (error) {
            utils.showAlert(`Error saving collection point: ${error.message}`, 'danger');
        }
    },
    
    deleteCollectionPoint: async (id) => {
        utils.displayConfirmationDialog('Are you sure you want to delete this collection point?', async () => {
            try {
                await api.deleteCollectionPoint(id);
                utils.showAlert('Collection point deleted successfully!', 'success');
                collectionPoints.displayCollectionPoints();
            } catch (error) {
                utils.showAlert(`Error deleting collection point: ${error.message}`, 'danger');
            }
        });
    }
};

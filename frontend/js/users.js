const users = {
    currentlyEditing: null,
    
    displayUsers: async () => {
        utils.clearContent();
        utils.setActiveNavItem('users-nav');
        
        const contentDiv = document.getElementById('content');
        
        try {
            const usersData = await api.getUsers();
            
            if (!Array.isArray(usersData)) {
                throw new Error('Unexpected response format');
            }
            
            contentDiv.innerHTML = `
                <div class="card">
                    <div class="card-header">
                        <h3>Users</h3>
                        <button class="btn btn-primary" id="add-user-btn">Add New User</button>
                    </div>
                    <div class="card-body">
                        <div class="table-responsive">
                            <table class="table table-striped">
                                <thead>
                                    <tr>
                                        <th>ID</th>
                                        <th>Name</th>
                                        <th>Email</th>
                                        <th>Actions</th>
                                    </tr>
                                </thead>
                                <tbody id="users-table-body">
                                    ${usersData.length > 0 ? usersData.map(user => `
                                        <tr>
                                            <td>${user.id}</td>
                                            <td>${user.name}</td>
                                            <td>${user.email}</td>
                                            <td>
                                                <button class="btn btn-sm btn-info btn-action edit-user" data-id="${user.id}">Edit</button>
                                                <button class="btn btn-sm btn-danger btn-action delete-user" data-id="${user.id}">Delete</button>
                                            </td>
                                        </tr>
                                    `).join('') : '<tr><td colspan="4" class="text-center">No users found</td></tr>'}
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
            `;
            
            // Add event listeners
            document.getElementById('add-user-btn').addEventListener('click', users.showAddUserModal);
            
            document.querySelectorAll('.edit-user').forEach(button => {
                button.addEventListener('click', (e) => {
                    const id = e.target.getAttribute('data-id');
                    users.showEditUserModal(id);
                });
            });
            
            document.querySelectorAll('.delete-user').forEach(button => {
                button.addEventListener('click', (e) => {
                    const id = e.target.getAttribute('data-id');
                    users.deleteUser(id);
                });
            });
            
        } catch (error) {
            contentDiv.innerHTML = `<div class="alert alert-danger">Error loading users: ${error.message}</div>`;
            console.error('Error in displayUsers:', error);
        }
    },
    
    showAddUserModal: () => {
        users.currentlyEditing = null;
        
        document.getElementById('modal-title').textContent = 'Add New User';
        
        const form = document.getElementById('item-form');
        form.innerHTML = `
            <input type="hidden" id="user-id" value="0">
            <div class="mb-3">
                <label for="user-name" class="form-label">Name</label>
                <input type="text" class="form-control" id="user-name" required>
            </div>
            <div class="mb-3">
                <label for="user-email" class="form-label">Email</label>
                <input type="email" class="form-control" id="user-email" required>
            </div>
            <div class="mb-3">
                <label for="user-password" class="form-label">Password</label>
                <input type="password" class="form-control" id="user-password" required>
            </div>
        `;
        
        document.getElementById('save-item').addEventListener('click', users.saveUser);
        
        const modal = new bootstrap.Modal(document.getElementById('itemModal'));
        modal.show();
    },
    
    showEditUserModal: async (id) => {
        try {
            const user = await api.getUserById(id);
            users.currentlyEditing = user;
            
            document.getElementById('modal-title').textContent = 'Edit User';
            
            const form = document.getElementById('item-form');
            form.innerHTML = `
                <input type="hidden" id="user-id" value="${user.id}">
                <div class="mb-3">
                    <label for="user-name" class="form-label">Name</label>
                    <input type="text" class="form-control" id="user-name" value="${user.name}" required>
                </div>
                <div class="mb-3">
                    <label for="user-email" class="form-label">Email</label>
                    <input type="email" class="form-control" id="user-email" value="${user.email}" required>
                </div>
                <div class="mb-3">
                    <label for="user-password" class="form-label">Password</label>
                    <input type="password" class="form-control" id="user-password" placeholder="Leave blank to keep current password">
                </div>
            `;
            
            document.getElementById('save-item').addEventListener('click', users.saveUser);
            
            const modal = new bootstrap.Modal(document.getElementById('itemModal'));
            modal.show();
            
        } catch (error) {
            utils.showAlert(`Error loading user: ${error.message}`, 'danger');
        }
    },
    
    saveUser: async () => {
        try {
            const id = document.getElementById('user-id').value;
            const name = document.getElementById('user-name').value;
            const email = document.getElementById('user-email').value;
            const password = document.getElementById('user-password').value;
            
            // Validate
            if (!name || !email) {
                utils.showAlert('Please fill all required fields', 'warning');
                return;
            }
            
            let userData = {
                id: parseInt(id),
                name: name,
                email: email
            };
            
            if (password) {
                userData.password = password;
            }
            
            let result;
            
            if (users.currentlyEditing) {
                result = await api.updateUser(id, userData);
                utils.showAlert('User updated successfully!', 'success');
            } else {
                if (!password) {
                    utils.showAlert('Password is required for new users', 'warning');
                    return;
                }
                result = await api.createUser(userData);
                utils.showAlert('User created successfully!', 'success');
            }
            
            const modal = utils.getModalInstance();
            modal.hide();
            
            // Refresh the users table
            users.displayUsers();
            
        } catch (error) {
            utils.showAlert(`Error saving user: ${error.message}`, 'danger');
        }
    },
    
    deleteUser: async (id) => {
        utils.displayConfirmationDialog('Are you sure you want to delete this user?', async () => {
            try {
                await api.deleteUser(id);
                utils.showAlert('User deleted successfully!', 'success');
                users.displayUsers();
            } catch (error) {
                utils.showAlert(`Error deleting user: ${error.message}`, 'danger');
            }
        });
    }
};

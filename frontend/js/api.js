const API_URL = 'http://localhost:1440/api';

const api = {
    // Reports
    getReports: async () => {
        try {
            console.log('Fetching reports...');
            const response = await fetch(`${API_URL}/reports`);
            
            console.log('Reports response status:', response.status);
            if (!response.ok) {
                throw new Error(`Failed to fetch reports: ${response.status}`);
            }
            
            const data = await response.json();
            console.log('Reports fetched successfully:', data);
            return data;
        } catch (error) {
            console.error('Error fetching reports:', error);
            throw error;
        }
    },
    
    getReportById: async (id) => {
        try {
            const response = await fetch(`${API_URL}/reports/${id}`);
            if (!response.ok) throw new Error(`Failed to fetch report with ID ${id}`);
            return await response.json();
        } catch (error) {
            console.error(`Error fetching report ${id}:`, error);
            throw error;
        }
    },
    
    createReport: async (reportData) => {
        try {
            const response = await fetch(`${API_URL}/reports`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(reportData)
            });
            if (!response.ok) throw new Error('Failed to create report');
            return await response.json();
        } catch (error) {
            console.error('Error creating report:', error);
            throw error;
        }
    },
    
    updateReport: async (id, reportData) => {
        try {
            const response = await fetch(`${API_URL}/reports/${id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(reportData)
            });
            if (!response.ok) throw new Error(`Failed to update report with ID ${id}`);
            return await response.json();
        } catch (error) {
            console.error(`Error updating report ${id}:`, error);
            throw error;
        }
    },
    
    deleteReport: async (id) => {
        try {
            const response = await fetch(`${API_URL}/reports/${id}`, {
                method: 'DELETE'
            });
            if (!response.ok) throw new Error(`Failed to delete report with ID ${id}`);
            return await response.json();
        } catch (error) {
            console.error(`Error deleting report ${id}:`, error);
            throw error;
        }
    },
    
    // Waste Types
    getWasteTypes: async () => {
        try {
            const response = await fetch(`${API_URL}/waste-types`);
            if (!response.ok) throw new Error('Failed to fetch waste types');
            return await response.json();
        } catch (error) {
            console.error('Error fetching waste types:', error);
            throw error;
        }
    },
    
    getWasteTypeById: async (id) => {
        try {
            const response = await fetch(`${API_URL}/waste-types/${id}`);
            if (!response.ok) throw new Error(`Failed to fetch waste type with ID ${id}`);
            return await response.json();
        } catch (error) {
            console.error(`Error fetching waste type ${id}:`, error);
            throw error;
        }
    },
    
    createWasteType: async (wasteTypeData) => {
        try {
            const response = await fetch(`${API_URL}/waste-types`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(wasteTypeData)
            });
            if (!response.ok) throw new Error('Failed to create waste type');
            return await response.json();
        } catch (error) {
            console.error('Error creating waste type:', error);
            throw error;
        }
    },
    
    updateWasteType: async (id, wasteTypeData) => {
        try {
            const response = await fetch(`${API_URL}/waste-types/${id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(wasteTypeData)
            });
            if (!response.ok) throw new Error(`Failed to update waste type with ID ${id}`);
            return await response.json();
        } catch (error) {
            console.error(`Error updating waste type ${id}:`, error);
            throw error;
        }
    },
    
    deleteWasteType: async (id) => {
        try {
            const response = await fetch(`${API_URL}/waste-types/${id}`, {
                method: 'DELETE'
            });
            if (!response.ok) throw new Error(`Failed to delete waste type with ID ${id}`);
            return await response.json();
        } catch (error) {
            console.error(`Error deleting waste type ${id}:`, error);
            throw error;
        }
    },
    
    // Collection Points
    getCollectionPoints: async () => {
        try {
            const response = await fetch(`${API_URL}/collection-points`);
            if (!response.ok) throw new Error('Failed to fetch collection points');
            return await response.json();
        } catch (error) {
            console.error('Error fetching collection points:', error);
            throw error;
        }
    },
    
    getCollectionPointById: async (id) => {
        try {
            const response = await fetch(`${API_URL}/collection-points/${id}`);
            if (!response.ok) throw new Error(`Failed to fetch collection point with ID ${id}`);
            return await response.json();
        } catch (error) {
            console.error(`Error fetching collection point ${id}:`, error);
            throw error;
        }
    },
    
    createCollectionPoint: async (collectionPointData) => {
        try {
            const response = await fetch(`${API_URL}/collection-points`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(collectionPointData)
            });
            if (!response.ok) throw new Error('Failed to create collection point');
            return await response.json();
        } catch (error) {
            console.error('Error creating collection point:', error);
            throw error;
        }
    },
    
    updateCollectionPoint: async (id, collectionPointData) => {
        try {
            const response = await fetch(`${API_URL}/collection-points/${id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(collectionPointData)
            });
            if (!response.ok) throw new Error(`Failed to update collection point with ID ${id}`);
            return await response.json();
        } catch (error) {
            console.error(`Error updating collection point ${id}:`, error);
            throw error;
        }
    },
    
    deleteCollectionPoint: async (id) => {
        try {
            const response = await fetch(`${API_URL}/collection-points/${id}`, {
                method: 'DELETE'
            });
            if (!response.ok) throw new Error(`Failed to delete collection point with ID ${id}`);
            return await response.json();
        } catch (error) {
            console.error(`Error deleting collection point ${id}:`, error);
            throw error;
        }
    },
    
    // Users
    getUsers: async () => {
        try {
            const response = await fetch(`${API_URL}/users`);
            if (!response.ok) throw new Error('Failed to fetch users');
            return await response.json();
        } catch (error) {
            console.error('Error fetching users:', error);
            throw error;
        }
    },
    
    getUserById: async (id) => {
        try {
            const response = await fetch(`${API_URL}/users/${id}`);
            if (!response.ok) throw new Error(`Failed to fetch user with ID ${id}`);
            return await response.json();
        } catch (error) {
            console.error(`Error fetching user ${id}:`, error);
            throw error;
        }
    },
    
    createUser: async (userData) => {
        try {
            const response = await fetch(`${API_URL}/users`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(userData)
            });
            if (!response.ok) throw new Error('Failed to create user');
            return await response.json();
        } catch (error) {
            console.error('Error creating user:', error);
            throw error;
        }
    },
    
    updateUser: async (id, userData) => {
        try {
            const response = await fetch(`${API_URL}/users/${id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(userData)
            });
            if (!response.ok) throw new Error(`Failed to update user with ID ${id}`);
            return await response.json();
        } catch (error) {
            console.error(`Error updating user ${id}:`, error);
            throw error;
        }
    },
    
    deleteUser: async (id) => {
        try {
            const response = await fetch(`${API_URL}/users/${id}`, {
                method: 'DELETE'
            });
            if (!response.ok) throw new Error(`Failed to delete user with ID ${id}`);
            return await response.json();
        } catch (error) {
            console.error(`Error deleting user ${id}:`, error);
            throw error;
        }
    }
};

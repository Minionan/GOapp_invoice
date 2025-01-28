// pages/clientEdit.js
// Fetch client details when page loads
document.addEventListener('DOMContentLoaded', () => {
    const urlParams = new URLSearchParams(window.location.search);
    const clientId = urlParams.get('id');

    if (!clientId || isNaN(clientId)) {
        alert('Invalid client ID');
        window.location.href = '/pages/client.html'; // Redirect to client list
        return;
    }

    // Fetch all clients
    fetch('/clients')
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to fetch clients');
            }
            return response.json();
        })
        .then(clients => {
            // Find the specific client by ID
            const client = clients.find(c => c.id === parseInt(clientId));
            if (!client) {
                throw new Error('Client not found');
            }

            // Populate the form fields
            document.getElementById('clientId').value = client.id;
            document.getElementById('clientName').value = client.clientName;
            document.getElementById('parentName').value = client.parentName;
            document.getElementById('address1').value = client.address1;
            document.getElementById('address2').value = client.address2;
            document.getElementById('phone').value = client.phone;
            document.getElementById('email').value = client.email;
            document.getElementById('abbreviation').value = client.abbreviation;
        })
        .catch(error => {
            console.error('Error fetching client details:', error);
            alert('Failed to fetch client details');
            window.location.href = '/pages/client.html'; // Redirect to client list
        });
});

function saveClient() {
    const clientId = document.getElementById('clientId').value;
    const clientName = document.getElementById('clientName').value;
    const parentName = document.getElementById('parentName').value;
    const address1 = document.getElementById('address1').value;
    const address2 = document.getElementById('address2').value;
    const phone = document.getElementById('phone').value;
    const email = document.getElementById('email').value;
    const abbreviation = document.getElementById('abbreviation').value;

    fetch(`/client-update?id=${clientId}`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ clientName, parentName, address1, address2, phone, email, abbreviation })
    })
    .then(response => {
        if (response.ok) {
            window.location.href = '/pages/client.html';
        } else {
            throw new Error('Failed to update client');
        }
    })
    .catch(error => {
        console.error('Error:', error);
        alert('Failed to update client');
    });
}

function deleteClient() {
    const clientId = document.getElementById('clientId').value;
    const confirmDelete = confirm('Are you sure you want to delete this client?');

    if (confirmDelete) {
        fetch(`/client-delete?id=${clientId}`, {
            method: 'POST'
        })
        .then(response => {
            if (response.ok) {
                window.location.href = '/pages/client.html';
            } else {
                throw new Error('Failed to delete client');
            }
        })
        .catch(error => {
            console.error('Error:', error);
            alert('Failed to delete client');
        });
    }
}
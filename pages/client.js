// pages/client.js
fetch('/clients')
// Fetch and display clients
.then(response => response.json())
.then(data => {
    //console.log('Received clients:', data); // debug line
    const tbody = document.querySelector('#clientTable tbody');
    tbody.innerHTML = '';
    data.forEach(client => {
        //console.log('Processing client:', client); // debug line
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${client.abbreviation}</td>
            <td>${client.clientName}</td>
            <td>${client.parentName}</td>
            <td>${client.phone}</td>
            <td>${client.email}</td>
            <td><button onclick="window.location.href='/pages/clientEdit.html?id=${client.id}'">Edit</button></td>
        `;
        tbody.appendChild(row);
    });
})
.catch(error => console.error('Error fetching clients:', error));

// Export clients to csv
function clientsExportCSV() {
    window.location.href = '/client-export';
}

// Clients importing from CSV
function importClientsCSV() {
    const fileInput = document.createElement('input');
    fileInput.type = 'file';
    fileInput.accept = '.csv';
    fileInput.onchange = (event) => {
        const file = event.target.files[0];
        if (!file) return;

        const formData = new FormData();
        formData.append('file', file);

        fetch('/client-import-csv', {
            method: 'POST',
            body: formData,
        })
        .then(response => response.json())
        .then(data => {
            let message = "Import Results:\n";
            message += `Imported Clients: ${(data.imported || []).join(', ') || 'None'}\n`;
            message += `Skipped Clients (duplicates): ${(data.skipped || []).join(', ') || 'None'}`;
            alert(message);
            window.location.href = '/pages/client.html'; // Reload client page
        })
        .catch(error => console.error('Error importing clients:', error));
    };

    fileInput.click();
}

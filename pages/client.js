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
